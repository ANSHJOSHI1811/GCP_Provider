package config
import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)
var DB *gorm.DB
var AuthToken = "Bearer ya29.a0AZYkNZhQ0SLHg9JEh9QU-5kDPFAB05qu_Xu0NkKCqAwSLQqfsmeja4SIJDmoFTF-QbDDYUD20uKDPAhZxfezJiUWCUBENjueHEKa_NfIVyocpsDgDZK0rSIWAkV4taNkFdVi7jOIH-NjA3IPlnRV2fkyku4uXsxag9YEJNk9ldw42FQaCgYKAYkSARMSFQHGX2MiccJNsobXfYVQhOq_2zIIFg0182" // your full token here
func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=password dbname=temp_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = database
}