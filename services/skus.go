package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gcp/config"
	"gcp/models"
)

type SkuResponse struct {
	Skus []SkuItem `json:"skus"`
}

type SkuItem struct {
	SkuID          string      `json:"skuId"`
	Category       SkuCategory `json:"category"`
	ServiceRegions []string    `json:"serviceRegions"`
}

type SkuCategory struct {
	ResourceFamily string `json:"resourceFamily"`
	UsageType      string `json:"usageType"`
}

func FetchAndInsertSkus() {
	url := "https://cloudbilling.googleapis.com/v1/services/6F81-5844-456A/skus"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}
	req.Header.Set("Authorization", config.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var skuResp SkuResponse
	if err := json.Unmarshal(body, &skuResp); err != nil {
		fmt.Println("❌ JSON unmarshal error:", err)
		fmt.Println("🧾 Raw body:", string(body))
		return
	}

	for _, sku := range skuResp.Skus {
		if len(sku.ServiceRegions) == 0 {
			continue
		}
		regionCode := sku.ServiceRegions[0]

		// Lookup region by region_code
		var region models.Region
		if err := config.DB.Where("region_code = ?", regionCode).First(&region).Error; err != nil {
			fmt.Printf("❌ Region not found in DB: %s\n", regionCode)
			continue
		}

		// Lookup provider using region.ProviderID
		var provider models.Provider
		if err := config.DB.Where("provider_id = ?", region.ProviderID).First(&provider).Error; err != nil {
			fmt.Printf("❌ Provider not found for region %s\n", regionCode)
			continue
		}

		// Check if SKU already exists
		var existing models.SKU
		if err := config.DB.Where("sku_code = ?", sku.SkuID).First(&existing).Error; err == nil {
			fmt.Printf("⚠️ SKU already exists: %s\n", sku.SkuID)
			continue
		}

		// Insert new SKU
		newSKU := models.SKU{
			ProviderID:    provider.ProviderID,
			RegionID:      region.RegionID,
			RegionCode:    region.RegionCode,
			SKUCode:       sku.SkuID,
			ProductFamily: sku.Category.ResourceFamily,
			Type:          sku.Category.UsageType,
			CreatedDate:   time.Now(),
			ModifiedDate:  time.Now(),
		}

		if err := config.DB.Create(&newSKU).Error; err != nil {
			fmt.Printf("❌ Failed to insert SKU %s: %v\n", sku.SkuID, err)
		} else {
			fmt.Printf("✅ Inserted SKU: %s\n", sku.SkuID)
		}
	}
}