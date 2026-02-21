# Tasks

## Sprint 1 - Core Pipeline

- [ ] Design notification schema and database tables
- [ ] Implement Redis-based priority queue
- [ ] Build notification router with preference resolution
- [ ] Create template engine with Go templates
- [ ] Implement SendGrid email channel
- [ ] Build in-app notification channel (WebSocket)
- [ ] Add delivery tracking per channel
- [ ] Create `POST /notifications` API endpoint
- [ ] Create `GET /notifications/:user_id` endpoint
- [ ] Write unit tests for template rendering

## Sprint 2 - Channels & Preferences

- [ ] Integrate Firebase Cloud Messaging for Android push
- [ ] Integrate APNs for iOS push
- [ ] Implement Twilio SMS channel
- [ ] Build device token registration API
- [ ] Create user preferences API (`GET` and `PUT`)
- [ ] Implement quiet hours with timezone support
- [ ] Add unsubscribe link generation for emails
- [ ] Write integration tests for each channel

## Sprint 3 - Batching & Analytics

- [ ] Implement 5-minute batch window aggregator
- [ ] Build daily/weekly digest scheduler
- [ ] Create digest summary template
- [ ] Add click tracking for email links
- [ ] Build notification analytics dashboard
- [ ] Performance test: 10k notifications/minute
- [ ] Document template authoring guide
- [ ] CAN-SPAM and GDPR compliance review
