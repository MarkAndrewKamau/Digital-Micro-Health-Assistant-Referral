package models

import (
	"time"

	"github.com/google/uuid"
)

type FacilityType string

const (
	FacilityTypeDispensary        FacilityType = "dispensary"
	FacilityTypeHealthCenter      FacilityType = "health_center"
	FacilityTypeSubCountyHospital FacilityType = "sub_county_hospital"
	FacilityTypeCountyHospital    FacilityType = "county_hospital"
	FacilityTypePrivateClinic     FacilityType = "private_clinic"
)

type Facility struct {
	ID               uuid.UUID              `json:"id"`
	Name             string                 `json:"name"`
	Type             FacilityType           `json:"type"`
	Level            *int                   `json:"level,omitempty"`
	County           *string                `json:"county,omitempty"`
	SubCounty        *string                `json:"sub_county,omitempty"`
	Latitude         *float64               `json:"latitude,omitempty"`
	Longitude        *float64               `json:"longitude,omitempty"`
	Address          *string                `json:"address,omitempty"`
	Phone            *string                `json:"phone,omitempty"`
	Email            *string                `json:"email,omitempty"`
	Services         []string               `json:"services"`
	OperatingHours   map[string]string      `json:"operating_hours"`
	AcceptsReferrals bool                   `json:"accepts_referrals"`
	AcceptsMpesa     bool                   `json:"accepts_mpesa"`
	BedCapacity      *int                   `json:"bed_capacity,omitempty"`
	StaffCount       *int                   `json:"staff_count,omitempty"`
	AvailableSlots   []map[string]interface{} `json:"available_slots"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type NearbyFacilitiesRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	RadiusKM  float64 `json:"radius_km" binding:"required,min=1,max=50"`
}