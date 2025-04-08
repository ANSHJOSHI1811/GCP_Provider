package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gcp/config"
	"gorm.io/gorm"
)

// API response struct
type APIRegion struct {
	Name string `json:"name"`
}

// Full API region list wrapper
type RegionList struct {
	Items []APIRegion `json:"items"`
}

// Provider DB model
type Provider struct {
	ProviderID   uint      `gorm:"primaryKey"`
	ProviderName string    `gorm:"unique"`
	CreatedDate  time.Time `gorm:"default:current_timestamp"`
	ModifiedDate time.Time `gorm:"default:current_timestamp"`
	DisableFlag  bool      `gorm:"default:false"`
}

// Region DB model
type Region struct {
	RegionID     uint      `gorm:"primaryKey"`
	RegionCode   string    `gorm:"unique"`
	ProviderID   uint      `gorm:"not null;constraint:OnDelete:CASCADE;"`
	CreatedDate  time.Time `gorm:"default:current_timestamp"`
	ModifiedDate time.Time `gorm:"default:current_timestamp"`
	DisableFlag  bool      `gorm:"default:false"`
}

func main() {
	// Connect to DB
	config.ConnectDatabase()

	// Step 1: Insert GCP provider if not exists
	var provider Provider
	err := config.DB.Where("provider_name = ?", "GCP").First(&provider).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		provider = Provider{
			ProviderName: "GCP",
			CreatedDate:  time.Now(),
			ModifiedDate: time.Now(),
			DisableFlag:  false,
		}
		if err := config.DB.Create(&provider).Error; err != nil {
			fmt.Println("❌ Failed to insert provider GCP:", err)
			return
		}
		fmt.Println("✅ Inserted new provider: GCP")
	} else if err != nil {
		fmt.Println("❌ Error checking provider:", err)
		return
	} else {
		fmt.Println("✅ Provider GCP already exists")
	}

	// Step 2: Call GCP API for regions
	url := "https://compute.googleapis.com/compute/v1/projects/812006687823/regions"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("❌ Request creation error:", err)
		return
	}
	req.Header.Set("Authorization", config.AuthToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("❌ Request error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("❌ Error reading response body:", err)
		return
	}

	var regionList RegionList
	err = json.Unmarshal(body, &regionList)
	if err != nil {
		fmt.Println("❌ JSON unmarshal error:", err)
		return
	}

	// Step 3: Insert regions into DB
	for _, apiRegion := range regionList.Items {
		var existing Region
		err := config.DB.Where("region_code = ?", apiRegion.Name).First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRegion := Region{
				RegionCode:   apiRegion.Name,
				ProviderID:   provider.ProviderID,
				CreatedDate:  time.Now(),
				ModifiedDate: time.Now(),
				DisableFlag:  false,
			}
			if err := config.DB.Create(&newRegion).Error; err != nil {
				fmt.Printf("❌ Failed to insert region %s: %v\n", apiRegion.Name, err)
			} else {
				fmt.Printf("✅ Inserted region: %s\n", apiRegion.Name)
			}
		} else if err != nil {
			fmt.Printf("❌ Error checking region %s: %v\n", apiRegion.Name, err)
		} else {
			fmt.Printf("⚠️ Region already exists: %s\n", apiRegion.Name)
		}
	}
}
