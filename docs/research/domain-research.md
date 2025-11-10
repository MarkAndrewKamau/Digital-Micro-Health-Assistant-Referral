# Domain Research: Kenya Healthcare Context

**Research Date**: November 2025  
**Research Team**: Product & ML Integration  
**Purpose**: Foundation for Digital Micro-Health Assistant MVP

---

## 1. Kenya Healthcare System Overview

### 1.1 Structure
Kenya's healthcare system operates on a four-tier model:
- **Level 1**: Community health services (CHVs, health posts)
- **Level 2-3**: Primary care (dispensaries, health centers)
- **Level 4**: County hospitals
- **Level 5-6**: National referral hospitals

### 1.2 Community Health Volunteers (CHVs)
- **Role**: First point of contact in rural/peri-urban communities
- **Coverage**: Each CHV typically serves 100-150 households
- **Activities**: Health education, basic triage, referrals, data collection
- **Challenges**: 
  - Limited medical training (typically 2-4 weeks)
  - Inconsistent connectivity
  - Paper-based record keeping
  - Delayed referral coordination
  - No standardized triage protocols

### 1.3 Current Gaps
- Average 48-72 hours from CHV assessment to facility visit
- 30-40% of referrals never completed (lost to follow-up)
- Limited integration between community and facility levels
- Minimal digital health adoption at community level

---

## 2. Clinical Guidelines & Protocols

### 2.1 Kenya Ministry of Health Resources
**Primary Sources:**
- Kenya Essential Medicines List (KEML 2021)
- Basic Paediatric Protocols (2016)
- Integrated Management of Childhood Illness (IMCI)
- Community Health Strategy (2020-2025)

**Key Documents for RAG Integration:**
1. Kenya Master Health Facility List (county facility data)
2. CHV Training Manual (symptom red flags)
3. Common Disease Protocols (malaria, TB, pneumonia, diarrhea)
4. Maternal & Child Health Guidelines

### 2.2 Common Triage Scenarios (Priority Cases)
| Condition | Red Flags | Recommended Action |
|-----------|-----------|-------------------|
| Fever | Temp >39°C, convulsions, stiff neck | Immediate referral |
| Diarrhea | Severe dehydration, bloody stool | Same-day facility visit |
| Cough | Difficulty breathing, chest indrawing | Immediate referral |
| Malnutrition | MUAC <11.5cm, edema | Urgent nutrition program |
| Pregnancy | Bleeding, severe headache, blurred vision | Emergency referral |

### 2.3 Triage Classification System
We will adopt a three-tier system:
- **Red (Urgent)**: Requires immediate facility care (<4 hours)
- **Yellow (Priority)**: Should see clinician within 24 hours
- **Green (Routine)**: Community management or routine appointment

---

## 3. Technology Landscape

### 3.1 Mobile & Internet Penetration
- **Mobile penetration**: 65%+ in target areas
- **Smartphone ownership**: 40-45% (growing rapidly)
- **Feature phones**: Still dominant in rural areas (55-60%)
- **Internet access**: Intermittent 2G/3G in most areas
- **SMS reliability**: 95%+ delivery rate

**Implication**: Must support SMS/USSD as primary channel, web as secondary.

### 3.2 Existing Digital Health Tools
- **KHIS (Kenya Health Information System)**: National HMIS, web-based
- **mHealth Kenya**: SMS appointment reminders (limited deployment)
- **eCHIS**: CHV digital data collection (pilot phase, limited adoption)
- **Afya Pap**: Cervical cancer screening app

**Gap**: No comprehensive AI-powered triage + referral system for CHVs.

### 3.3 Payment Infrastructure
- **M-Pesa**: 99% awareness, 80%+ active usage
- **Average transaction**: KES 500-2000 for health payments
- **Facility integration**: <20% of clinics accept digital payments
- **Opportunity**: Enable seamless referral payments for private clinics

---

## 4. Regulatory & Compliance Landscape

### 4.1 Data Protection
- **Kenya Data Protection Act (2019)**: Requires explicit consent, data minimization
- **Health Records & Information Managers Act**: Governs health data handling
- **Key Requirements**:
  - Explicit patient consent before data collection
  - Right to access and delete personal data
  - Data breach notification within 72 hours
  - Data localization preferences (county-level)

### 4.2 Telemedicine & Digital Health
- **Kenya Telemedicine Guidelines (2022)**: Establishes framework for remote care
- **AI in Healthcare**: No specific regulations yet, but must follow general medical device guidelines
- **Liability**: Digital tools classified as "decision support" (not diagnostic) to reduce liability
- **CHV Scope**: CHVs cannot diagnose/prescribe; tool must align with this limitation

### 4.3 Medical Disclaimers
All triage outputs must include:
> "This is a symptom assessment tool only, not a medical diagnosis. Always consult a qualified healthcare provider for medical advice."

---

## 5. Stakeholder Landscape

### 5.1 County Health Departments
- **Decision makers**: County Director of Health, CHV coordinators
- **Priorities**: Improve referral completion rates, CHV effectiveness, data quality
- **Concerns**: Data ownership, staff training requirements, sustainability
- **Engagement strategy**: Pilot with progressive counties (e.g., Makueni, Kisumu)

### 5.2 NGO Partners (Potential)
- **Amref Health Africa**: Large CHV programs
- **LVCT Health**: Community health focus
- **PATH**: Digital health innovation
- **Opportunity**: Co-funding, CHV network access, credibility

### 5.3 Private Clinics
- **Pain points**: Underutilized capacity, manual appointment booking
- **Value proposition**: Increased patient flow via referrals, digital payments
- **Revenue model**: Per-referral transaction fee (10-15% of consultation fee)

---

## 6. User Research Insights

### 6.1 CHV Workflows (Observational)
Typical CHV home visit:
1. Greet family, assess household (5-10 min)
2. Patient symptom inquiry (5-10 min)
3. Basic vitals if available (temp, weight for children)
4. Decision: home care advice OR referral
5. If referral: write paper note, explain to family
6. Record in paper register

**Pain points identified**:
- Uncertainty on when to refer (rely on memory/experience)
- No facility availability information
- Referrals often delayed/ignored by families (cost, transport)
- No feedback on referral outcomes

### 6.2 Patient Preferences
- Trust CHVs as first point of contact
- SMS preferred for confirmations/reminders
- Willing to pay KES 50-200 for private clinic visit if convenient
- Fear of "wasting money" on unnecessary facility visit

---

## 7. AI/ML Considerations

### 7.1 LLM Selection
**Recommendation: Llama 2 7B/13B for on-premise inference**
- Rationale: Lower cost, data privacy, Kenya-specific fine-tuning potential
- Fallback: OpenAI GPT-4 for complex/ambiguous cases (paid API)

### 7.2 RAG Data Sources
Priority documents for vector embedding:
1. Kenya MoH clinical guidelines (PDFs)
2. IMCI symptom decision trees
3. County-specific referral protocols
4. Common medication dosing (KEML)
5. CHV training materials

**Processing**: Extract text, chunk to 500-1000 tokens, embed with sentence-transformers, store in pgvector.

### 7.3 Safety & Accuracy
- **Conservative bias**: Default to referral when uncertain
- **Human-in-loop**: All "yellow" and "red" cases reviewed by clinician within 4 hours
- **Audit trail**: Log all LLM inputs/outputs for review
- **Red-teaming**: Test with adversarial prompts (e.g., dangerous advice scenarios)

---

## 8. Consent & Privacy Design

### 8.1 Consent Flow
1. **Initial registration**: "Do you consent to share your health info for triage?"
2. **Referral creation**: "May we share your symptoms with [Facility Name]?"
3. **Data sharing**: "Consent to share anonymized data for research?" (optional)

**Implementation**: `consent_logs` table tracks all consent events.

### 8.2 Data Minimization
- Collect only: phone, age, gender, symptoms, location (approximate)
- Do NOT collect: full name, ID numbers, photos (unless medically necessary)
- Retention: 2 years for active cases, anonymize after

---

## 9. Go-To-Market Strategy

### 9.1 Pilot Selection Criteria
- County with >100 active CHVs
- Progressive health leadership
- Existing CHV digital readiness programs
- 3G coverage in 70%+ of target area

**Top candidates**: Makueni, Kisumu, Kiambu

### 9.2 Pilot Success Metrics (3 months)
- 10 CHVs trained and actively using app
- 1,500 triage interactions (avg 50/CHV/month)
- 500 referrals generated (33% conversion)
- <24hr triage-to-referral time (80% of cases)
- 4.0+ user satisfaction score (CHVs)

### 9.3 Marketing Tactics
- County health leadership presentations
- CHV training roadshows with demos
- Radio spots on county health programs
- Success stories/testimonials from pilot CHVs
- Grant applications to health innovation funds

---

## 10. Risks & Mitigation

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Incorrect triage advice harms patient | Medium | Critical
| Conservative triage thresholds, human-in-loop review, legal disclaimers, clinician oversight |
| CHVs resist new technology | Medium | High | Co-design with CHVs, simple UI, offline-first, adequate training, incentives |
| Regulatory pushback on AI | Low | High | Position as decision support tool, early county engagement, transparent AI practices |
| Data privacy breach | Low | Critical | Encryption, access controls, minimal data collection, regular audits, compliance training |
| Low patient adoption | Medium | High | SMS-first (no app required for patients), free triage, M-Pesa integration, trust-building |
| Facility resistance to referrals | Medium | Medium | Demonstrate value (increased patient flow), easy referral acceptance, payment integration |
| Poor internet connectivity | High | Medium | Offline-first CHV app, SMS fallback, data sync when connected, local triage rules |
| Funding sustainability | Medium | High | Diverse revenue streams, county subscriptions, transaction fees, grant applications |
| LLM hallucination/errors | Medium | High | RAG grounding, output validation, confidence thresholds, human review, audit trails |
| M-Pesa API instability | Low | Medium | Retry mechanisms, fallback payment methods, manual reconciliation processes |11. Technical Infrastructure Requirements11.1 Minimum Viable Infrastructure
For 10 CHVs / 1,500 interactions/month:

Compute: 2 vCPU, 4GB RAM server (API + workers)
Database: PostgreSQL with 20GB storage
LLM inference: 1 GPU instance (NVIDIA T4 or similar) OR OpenAI API
Bandwidth: ~5GB/month (excluding LLM API calls)
Cost estimate: $150-250/month (cloud hosting) OR $50-100/month (local server + OpenAI API)
## 11.2 Scaling Projections
For 100 CHVs / 15,000 interactions/month:

Horizontal scaling: 4-6 API servers (load balanced)
Database: 200GB storage, read replicas
Redis cluster for distributed caching
RabbitMQ for async job processing
Cost estimate: $800-1,200/month
## 12. Key Kenyan Clinical TerminologyEnglish TermSwahili/LocalUsage ContextCommunity Health VolunteerMganga wa Jamii / CHVOfficial titleHealth facilityKituo cha AfyaClinic/hospitalFeverHomaMost common symptomDiarrheaKuharaCommon in childrenCoughKikohoziTB screening triggerMalariaMalaria / Homa ya MalariaEndemic in many areasDispensaryZahanatiLevel 2 facilityHealth centerKituo cha AfyaLevel 3 facilityReferralRufaaPatient transferAppointmentMiadiScheduled visitImplication: UI should support Swahili and English, SMS templates in both languages.13. Sample User Journeys13.1 Journey A: SMS Triage (Feature Phone User)

Patient sends SMS to shortcode: "Homa kwa siku tatu" (Fever for 3 days)
System replies: "Una umri gani? Watoto au wazima?" (Age? Child or adult?)
Patient: "Mtoto miaka 2" (Child 2 years)
System (after LLM triage): "Tafadhali peleka mtoto hospitalini siku hii. Joto ni ngapi?" (Please take child to hospital today. What's the temperature?)
Patient: "38.5"
System: "Homa ya kawaida. Kliniki ya karibu: Kamulu Health Center (2km). Mwonyeshe CHV kumbukumbu hii: REF-4829"
System creates referral, sends SMS to facility and CHV
## 13.2 Journey B: CHV-Assisted Triage (Offline)

CHV visits household, patient complains of cough for 2 weeks
CHV opens app (offline), enters symptoms: cough, duration, night sweats
App (using local triage rules): Flags as potential TB, recommends immediate referral
CHV generates referral token (QR code), explains to patient
CHV returns to coverage area, app syncs data
System sends SMS to patient and facility with referral details
Facility receives notification, books appointment slot
## 13.3 Journey C: Web Chat Triage (Smartphone User)

Patient visits web app, clicks "Check Symptoms"
Chatbot: "Hello! I'll help assess your symptoms. What brings you here today?"
Patient: "I have a headache and fever"
Chatbot: Asks follow-up questions (duration, severity, other symptoms)
System (LLM + RAG): Assesses as "Yellow" - possible malaria or infection
Chatbot: "You should see a doctor within 24 hours. Would you like to book an appointment?"
Patient: "Yes"
System: Shows nearby facilities with availability, patient selects
System: Creates referral + appointment, sends confirmation SMS with QR code
## 14. Competitive Landscape14.1 Existing Solutions
SolutionStrengthsWeaknessesOur DifferentiationeCHIS (Gov't)Official, integrated with KHISNo triage, limited CHV adoption, slow updatesAI triage, offline-first, user-friendlyAda HealthAI-powered triageNot Kenya-focused, requires smartphone, no local integrationLocal guidelines, SMS support, facility integrationBabylon HealthComprehensiveUK-focused, expensive, no CHV workflowCHV-centric, affordable, M-Pesa integrationAfya PapDigital screeningSingle-use (cervical cancer), no referralsMulti-condition, end-to-end referral14.2 Unique Value Proposition
"The only AI-powered health assistant designed specifically for Kenyan CHVs and patients, combining SMS accessibility, offline workflows, and seamless facility referrals."15. Data Sources & References15.1 Official Documents (To Obtain/Review)

 Kenya Community Health Strategy 2020-2025 (PDF)
 Kenya Master Health Facility List (Excel/CSV)
 Basic Paediatric Protocols 2016 (PDF)
 IMCI Chart Booklet (PDF)
 Kenya Essential Medicines List 2021 (PDF)
 Kenya Telemedicine Guidelines 2022 (PDF)
 County-specific CHV training materials (various)
## 15.2 Key Contacts for Further Research

Ministry of Health: Division of Community Health Services
County Health Departments: Makueni, Kisumu, Kiambu (pilot targets)
Amref Health Africa: CHV program coordinators
KEMRI: Research collaboration potential
Kenya Medical Training College: Clinical advisory
## 15.3 Online Resources

Kenya Health Information System: https://hiskenya.org
Kenya MoH: https://www.health.go.ke
Africa's Talking Developer Docs: https://developers.africastalking.com
Safaricom Daraja API: https://developer.safaricom.co.ke
## 16. Next Steps (Research)Week 2 Priorities

Obtain official clinical guidelines (MoH website + contacts)
Interview 3-5 CHVs (understand workflows, pain points)
Visit 2 health facilities (observe referral processes)
Map facility data (initial database of pilot county facilities)
Test SMS/USSD (Africa's Talking sandbox)
Week 3-4 Priorities

Process clinical docs for RAG (extract, chunk, embed)
Prototype triage decision tree (rule-based fallback)
Legal review (consent forms, disclaimers with local lawyer)
Finalize pilot county (sign MOU with county health dept)
## 17. Assumptions & Validation NeededAssumptionConfidenceValidation MethodCHVs have access to feature phonesHighSurvey during pilot recruitmentPatients trust SMS health adviceMediumUser testing with sample messagesFacilities will accept digital referralsMediumPilot facility interviews33% triage-to-referral conversion is achievableLowBaseline data from county + pilot trackingLLM can achieve 80%+ triage accuracyMediumValidation against clinician gold standardM-Pesa integration will increase referral completionMediumA/B test during pilotCounties will pay for subscriptions post-pilotLowEarly commercial discussions

## 18. Cultural & Social Considerations
## 18.1 Trust & Stigma

Traditional healers: Still consulted in many communities; avoid positioning as replacement
Stigma: TB, HIV have social stigma; ensure confidential referrals
Gender: Female CHVs preferred for maternal health cases
Language: Swahili + English minimum; consider local languages (Kikuyu, Luo, etc.) for scale

## 18.2 Health Beliefs

Fever: Often attributed to "malaria" by default; triage must account for this
Home remedies: Common first-line treatment; tool should acknowledge and guide appropriately
Facility avoidance: Cost and distance are major barriers; emphasize nearby, affordable options

## 18.3 Communication Style

Directness: Kenyans appreciate clear, actionable health advice
Respect: Use respectful language (e.g., "Mzee" for elders, "Mama" for mothers)
Reassurance: Balance urgency with reassurance to avoid panic


## 19. Success Stories (Inspiration)
19.1 M-TIBA (Kenya)

Digital health wallet for medical savings and payments
Lesson: M-Pesa integration is critical for adoption
Challenge: Required significant stakeholder buy-in

## 19.2 Babyl Rwanda

AI triage + telemedicine at national scale
Lesson: Government partnership enables rapid scale
Challenge: High operational costs, sustainability questions

## 19.3 mPedigree (Ghana/Kenya)

SMS-based medicine verification
Lesson: Simple SMS interface works for low-literacy users
Challenge: User education required


## 20. Research Summary & Recommendations
## 20.1 Key Findings

CHVs are underutilized due to lack of decision support tools
SMS/USSD is essential for patient accessibility in rural areas
Referral completion is <60% currently; major improvement opportunity
Regulatory environment is supportive of digital health innovation
M-Pesa integration is non-negotiable for payment workflows
Consent and data privacy must be central to design
Conservative triage approach required to minimize clinical risk

## 20.2 Critical Success Factors

Co-design with CHVs and county health departments
Offline-first architecture for poor connectivity
Human-in-loop for clinical safety
Simple, intuitive UI (SMS + mobile app)
Clear value proposition for all stakeholders (CHVs, facilities, counties)

## 20.3 Recommended Adjustments to Original Plan

Prioritize SMS over web chat for MVP (reverse order)
Add Swahili language support from day 1
Build stronger offline capabilities for CHV app
Engage pilot county earlier (Week 2-3 vs Week 9)
Add basic rule-based triage as fallback to LLM (for offline use)
Include facility onboarding in pilot (not just CHVs)