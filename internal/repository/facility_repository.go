package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/MarkAndrewKamau/Digital-Micro-Health-Assistant-Referral/internal/models"
)

type FacilityRepository struct {
	db *pgxpool.Pool
}

func NewFacilityRepository(db *pgxpool.Pool) *FacilityRepository {
	return &FacilityRepository{db: db}
}

func (r *FacilityRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Facility, error) {
	query := `
		SELECT 
			id, name, type, level, county, sub_county,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			address, phone, email, services, operating_hours,
			accepts_referrals, accepts_mpesa, bed_capacity, staff_count,
			available_slots, created_at, updated_at
		FROM facilities
		WHERE id = $1
	`

	var facility models.Facility
	var servicesRaw, operatingHoursRaw, availableSlotsRaw []byte

	err := r.db.QueryRow(ctx, query, id).Scan(
		&facility.ID,
		&facility.Name,
		&facility.Type,
		&facility.Level,
		&facility.County,
		&facility.SubCounty,
		&facility.Latitude,
		&facility.Longitude,
		&facility.Address,
		&facility.Phone,
		&facility.Email,
		&servicesRaw,
		&operatingHoursRaw,
		&facility.AcceptsReferrals,
		&facility.AcceptsMpesa,
		&facility.BedCapacity,
		&facility.StaffCount,
		&availableSlotsRaw,
		&facility.CreatedAt,
		&facility.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("facility not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get facility: %w", err)
	}

	if err := json.Unmarshal(servicesRaw, &facility.Services); err != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", err)
	}
	if err := json.Unmarshal(operatingHoursRaw, &facility.OperatingHours); err != nil {
		return nil, fmt.Errorf("failed to unmarshal operating hours: %w", err)
	}
	if err := json.Unmarshal(availableSlotsRaw, &facility.AvailableSlots); err != nil {
		return nil, fmt.Errorf("failed to unmarshal available slots: %w", err)
	}

	return &facility, nil
}

func (r *FacilityRepository) GetNearby(ctx context.Context, lat, lng, radiusKM float64) ([]*models.Facility, error) {
	query := `
		SELECT 
			id, name, type, level, county, sub_county,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			address, phone, email, services, operating_hours,
			accepts_referrals, accepts_mpesa, bed_capacity, staff_count,
			available_slots, created_at, updated_at,
			ST_Distance(location, ST_GeogFromText($1)) / 1000 as distance_km
		FROM facilities
		WHERE 
			accepts_referrals = true
			AND ST_DWithin(location, ST_GeogFromText($1), $2)
		ORDER BY distance_km ASC
		LIMIT 20
	`

	point := fmt.Sprintf("POINT(%f %f)", lng, lat)
	radiusMeters := radiusKM * 1000

	rows, err := r.db.Query(ctx, query, point, radiusMeters)
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby facilities: %w", err)
	}
	defer rows.Close()

	var facilities []*models.Facility
	for rows.Next() {
		var facility models.Facility
		var servicesRaw, operatingHoursRaw, availableSlotsRaw []byte
		var distanceKM float64

		err := rows.Scan(
			&facility.ID,
			&facility.Name,
			&facility.Type,
			&facility.Level,
			&facility.County,
			&facility.SubCounty,
			&facility.Latitude,
			&facility.Longitude,
			&facility.Address,
			&facility.Phone,
			&facility.Email,
			&servicesRaw,
			&operatingHoursRaw,
			&facility.AcceptsReferrals,
			&facility.AcceptsMpesa,
			&facility.BedCapacity,
			&facility.StaffCount,
			&availableSlotsRaw,
			&facility.CreatedAt,
			&facility.UpdatedAt,
			&distanceKM,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan facility: %w", err)
		}

		if err := json.Unmarshal(servicesRaw, &facility.Services); err != nil {
			return nil, fmt.Errorf("failed to unmarshal services: %w", err)
		}
		if err := json.Unmarshal(operatingHoursRaw, &facility.OperatingHours); err != nil {
			return nil, fmt.Errorf("failed to unmarshal operating hours: %w", err)
		}
		if err := json.Unmarshal(availableSlotsRaw, &facility.AvailableSlots); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available slots: %w", err)
		}

		facilities = append(facilities, &facility)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating facilities: %w", err)
	}

	return facilities, nil
}

func (r *FacilityRepository) List(ctx context.Context, county *string, facilityType *models.FacilityType) ([]*models.Facility, error) {
	query := `
		SELECT 
			id, name, type, level, county, sub_county,
			ST_Y(location::geometry) as latitude,
			ST_X(location::geometry) as longitude,
			address, phone, email, services, operating_hours,
			accepts_referrals, accepts_mpesa, bed_capacity, staff_count,
			available_slots, created_at, updated_at
		FROM facilities
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if county != nil {
		query += fmt.Sprintf(" AND county = $%d", argCount)
		args = append(args, *county)
		argCount++
	}

	if facilityType != nil {
		query += fmt.Sprintf(" AND type = $%d", argCount)
		args = append(args, *facilityType)
		argCount++
	}

	query += " ORDER BY name ASC LIMIT 100"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list facilities: %w", err)
	}
	defer rows.Close()

	var facilities []*models.Facility
	for rows.Next() {
		var facility models.Facility
		var servicesRaw, operatingHoursRaw, availableSlotsRaw []byte

		err := rows.Scan(
			&facility.ID,
			&facility.Name,
			&facility.Type,
			&facility.Level,
			&facility.County,
			&facility.SubCounty,
			&facility.Latitude,
			&facility.Longitude,
			&facility.Address,
			&facility.Phone,
			&facility.Email,
			&servicesRaw,
			&operatingHoursRaw,
			&facility.AcceptsReferrals,
			&facility.AcceptsMpesa,
			&facility.BedCapacity,
			&facility.StaffCount,
			&availableSlotsRaw,
			&facility.CreatedAt,
			&facility.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan facility: %w", err)
		}

		if err := json.Unmarshal(servicesRaw, &facility.Services); err != nil {
			return nil, fmt.Errorf("failed to unmarshal services: %w", err)
		}
		if err := json.Unmarshal(operatingHoursRaw, &facility.OperatingHours); err != nil {
			return nil, fmt.Errorf("failed to unmarshal operating hours: %w", err)
		}
		if err := json.Unmarshal(availableSlotsRaw, &facility.AvailableSlots); err != nil {
			return nil, fmt.Errorf("failed to unmarshal available slots: %w", err)
		}

		facilities = append(facilities, &facility)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating facilities: %w", err)
	}

	return facilities, nil
}