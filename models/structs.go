package models

import "time"

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

type SKU struct {
	ID                   uint      `gorm:"primaryKey"`
	RegionID             uint      `gorm:"not null;constraint:OnDelete:CASCADE;"` // Foreign key with cascade delete
	ProviderID           uint      `gorm:"not null"`
	RegionCode           string    `gorm:"not null"`
	SKUCode              string    `gorm:"unique"`
	ArmSkuName           string    `gorm:"column:arm_sku_name"`
	InstanceSKU          string
	ProductFamily        string
	VCPU                 int
	CpuArchitecture      string
	InstanceType         string    `gorm:"column:instance_type"` // Storing 'name' in instance_type
	Storage              string
	Network              string
	OperatingSystem      string
	Type                 string
	Memory               string
	PhysicalProcessor    string    `gorm:"column:physical_processor"`
	MaxThroughput        string    `gorm:"column:max_throughput"`
	EnhancedNetworking   string    `gorm:"column:enhanced_networking"`
	GPU                  string    `gorm:"column:gpu"`
	MaxIOPS              string    `gorm:"column:max_iops"`
	CreatedDate          time.Time `gorm:"default:current_timestamp"`
	ModifiedDate         time.Time `gorm:"default:current_timestamp"`
	DisableFlag          bool      `gorm:"default:false"`
}

// In sku.go (API logic)
type SkuResponse struct {
    Skus []SkuItem `json:"skus"`
}

type SkuItem struct {
    SkuID          string      `json:"skuId"`
    Name           string      `json:"name"`
    Category       SkuCategory `json:"category"`
    ServiceRegions []string    `json:"serviceRegions"`
}

type SkuCategory struct {
    ResourceFamily string `json:"resourceFamily"`
}
