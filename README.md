# Digital Micro-Health Assistant & Referral Platform

## Overview
AI-enabled triage assistant for rural and peri-urban patients in Kenya, connecting Community Health Volunteers (CHVs), patients, and healthcare facilities through SMS/USSD and web interfaces.

## Target Users
- Patients in rural/peri-urban areas
- Community Health Volunteers (CHVs)
- County health facilities and clinics
- County health departments
- NGO partners

## Value Proposition
Provide accessible healthcare triage through:
- AI-powered symptom assessment via SMS/USSD and web chat
- Automated appointment scheduling and referrals
- Offline-first CHV workflows
- Integration with Kenya's health infrastructure

## KPI Targets (MVP)
- **10 pilot CHVs** deployed and trained
- **1,500 patient interactions** processed
- **500 successful referrals** completed
- **<24 hours** average triage-to-referral time

## Core Features (MVP)
1. **Symptom Triage** - LLM-powered assessment via SMS/USSD and web
2. **Appointment Booking** - Facility scheduling with QR-based referral tokens
3. **CHV Mobile App** - Offline data capture with sync capability
4. **Clinician Escalation** - Human-in-loop for complex cases
5. **Payment Integration** - M-Pesa for private clinic referrals
6. **Consent Management** - Privacy-first data handling

## Technology Stack
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL with PostGIS (geospatial) and pgvector (AI embeddings)
- **Cache**: Redis
- **Message Queue**: RabbitMQ
- **Storage**: MinIO (S3-compatible)
- **AI**: Llama-2 (on-prem) / OpenAI (cloud fallback)
- **SMS/USSD**: Africa's Talking / Twilio
- **Payments**: Safaricom Daraja (M-Pesa)
- **Observability**: OpenTelemetry, Prometheus, Grafana, Sentry

## Architecture Services
- **API Gateway** - Rate limiting and routing
- **Auth & Consent Service** - Identity and consent management
- **Triage Service** - LLM + RAG symptom assessment
- **Scheduling Service** - Appointments and referral tokens
- **Notifications Service** - SMS/USSD/Email delivery
- **Clinic Portal** - Web interface for facility staff
- **Escalation Worker** - Route complex cases to clinicians

## Database Schema (Core Tables)
- `patients` - Patient demographics and consent
- `triage_sessions` - Symptom assessment records
- `referrals` - Referral tracking with tokens
- `facilities` - Healthcare facility data with geolocation
- `clinicians` - Healthcare provider information
- `appointments` - Scheduled patient visits
- `embeddings` - Vector store for RAG
- `consent_logs` - Audit trail for consent actions

## Development Setup

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- Git

### Quick Start
```bash
# Clone repository
git clone <repo-url>
cd digital-health-assistant

# Copy environment template
cp .env.example .env

# Start infrastructure services
docker-compose up -d

# Install dependencies
go mod download

# Run migrations (once available)
# make migrate-up

# Start development server
go run cmd/server/main.go
```

Server will be available at `http://localhost:8080`

### Health Check
```bash
curl http://localhost:8080/health
```

## Environment Variables
See `.env.example` for required configuration. Key variables:
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `RABBITMQ_URL` - RabbitMQ connection string
- `AFRICASTALKING_*` - SMS/USSD credentials
- `MPESA_*` - M-Pesa Daraja API credentials
- `OPENAI_API_KEY` - Optional LLM fallback

## Project Epics
1. **Auth & Consent** - User identity and privacy controls
2. **Triage** - Symptom assessment engine
3. **Referrals** - Facility coordination and tokens
4. **CHV App** - Mobile offline-first application
5. **Payments** - M-Pesa integration
6. **Integrations** - SMS/USSD, Maps, External APIs

## MVP Timeline (12 weeks)
- **Weeks 1-4**: Core API, SMS/USSD, basic triage, CHV app skeleton
- **Weeks 5-8**: LLM/RAG integration, clinician escalation
- **Weeks 9-12**: Scheduling, referral tokens, pilot deployment

## Security & Compliance
- Explicit consent capture before PII collection
- HIPAA-aligned best practices
- Encrypted data at rest and in transit
- Role-based access control (RBAC)
- Audit trails for all clinical actions
- Regular LLM output review

## Monetization Strategy
- County health department subscriptions
- Per-referral transaction fees (private clinics)
- Premium AI triage features
- NGO/grant partnerships

## Contributing
[Guidelines to be added]

## License


## Contact
gmail: kamaumark19@gmail.com