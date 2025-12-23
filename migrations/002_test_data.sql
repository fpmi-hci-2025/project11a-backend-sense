-- Test data for handler testing
-- All UUIDs are fixed for consistency in testing
-- Password hashes are bcrypt with DefaultCost

BEGIN;

-- USERS (7 users with different roles)
INSERT INTO users (id, username, email, phone, icon_url, description, role, password_hash, registered_at, followers_count, following_count, statistic) VALUES
-- Super admin
('00000000-0000-0000-0000-000000000001', 'superadmin', 'super@example.com', '+375291111111', 'https://example.com/icons/super.png', 'Super administrator account', 'super', '$2a$10$/htG7FOroStzjQ4k5PTdBOQEhG9zpjd9q/4TksD2GrZ0k18gia.ze', '2024-01-01 10:00:00+00', 5, 2, '{"publications": 3, "likes": 15, "comments": 8}'),
-- Expert
('00000000-0000-0000-0000-000000000002', 'expert_user', 'expert@example.com', '+375292222222', 'https://example.com/icons/expert.png', 'Expert content creator with deep knowledge', 'expert', '$2a$10$1djuS5xytqkgzth7FfykleNMCOTr8SOf76xRIkkMn9f6zIcUxMj/m', '2024-01-15 12:30:00+00', 120, 45, '{"publications": 25, "likes": 350, "comments": 89}'),
-- Creator
('00000000-0000-0000-0000-000000000003', 'content_creator', 'creator@example.com', '+375293333333', 'https://example.com/icons/creator.png', 'Creative content maker', 'creator', '$2a$10$Oxa3IHXMj8kHQ0hA1FaeH.Hrw/ZQaWiV.u0FYMtxXEFNI..8mhdSG', '2024-02-01 09:15:00+00', 89, 67, '{"publications": 18, "likes": 234, "comments": 56}'),
-- Regular user 1
('00000000-0000-0000-0000-000000000004', 'testuser1', 'testuser1@example.com', '+375294444444', 'https://example.com/icons/user1.png', 'First test user for API testing', 'user', '$2a$10$yHrCLvWjsEs4cQf7fYjzmOHfb34i7IA/bLawLEzwWZcgZUUXq5Qe.', '2024-02-10 14:20:00+00', 23, 12, '{"publications": 5, "likes": 45, "comments": 12}'),
-- Regular user 2
('00000000-0000-0000-0000-000000000005', 'testuser2', 'testuser2@example.com', '+375295555555', 'https://example.com/icons/user2.png', 'Second test user for API testing', 'user', '$2a$10$mbM/so1uXWTeUOUkH16zLeKiarW/2KOEZtP9oIXq5l1me7WheCdgG', '2024-02-12 16:45:00+00', 15, 8, '{"publications": 3, "likes": 28, "comments": 7}'),
-- Regular user 3
('00000000-0000-0000-0000-000000000006', 'regular_user', 'regular@example.com', '+375296666666', NULL, 'Regular user without icon', 'user', '$2a$10$qY8RzpnOQplFiQc4a/CVde84opreD0Zk/NOjweaYN0PcdM/utuG/i', '2024-03-01 11:00:00+00', 7, 5, '{"publications": 2, "likes": 12, "comments": 3}'),
-- Reader
('00000000-0000-0000-0000-000000000007', 'reader_only', 'reader@example.com', NULL, NULL, 'Reader account with no publications', 'reader', '$2a$10$yHrCLvWjsEs4cQf7fYjzmOHfb34i7IA/bLawLEzwWZcgZUUXq5Qe.', '2024-03-15 08:30:00+00', 0, 3, '{"publications": 0, "likes": 5, "comments": 2}');

-- PUBLICATIONS (20 publications of different types)
INSERT INTO publications (id, author_id, type, title, content, source, publication_date, visibility) VALUES
-- Quotes
('10000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', 'quote', 'The only way to do great work is to love what you do.', 'The only way to do great work is to love what you do.', 'Steve Jobs', '2024-01-20 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003', 'quote', 'Innovation distinguishes between a leader and a follower.', 'Innovation distinguishes between a leader and a follower.', 'Steve Jobs', '2024-01-25 14:30:00+00', 'public'),
('10000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000001', 'quote', 'Stay hungry, stay foolish.', 'Stay hungry, stay foolish.', 'Steve Jobs', '2024-02-01 09:15:00+00', 'community'),
('10000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000004', 'quote', 'Life is what happens to you while you''re busy making other plans.', 'Life is what happens to you while you''re busy making other plans.', 'John Lennon', '2024-02-15 16:20:00+00', 'public'),
('10000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000005', 'quote', 'The future belongs to those who believe in the beauty of their dreams.', 'The future belongs to those who believe in the beauty of their dreams.', 'Eleanor Roosevelt', '2024-02-20 11:45:00+00', 'private'),

-- Posts
('10000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000002', 'post', 'Amazing Book on Software Architecture', 'Just finished reading an amazing book about software architecture. The principles of clean code are more important than ever in today''s fast-paced development environment.', NULL, '2024-02-05 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000003', 'post', 'New AI + Web Development Project', 'Working on a new project that combines AI with traditional web development. Excited to share more details soon!', NULL, '2024-02-10 10:30:00+00', 'public'),
('10000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000004', 'post', 'Database Optimization Learning', 'Today I learned something new about database optimization. Always keep learning!', NULL, '2024-02-18 15:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000005', 'post', 'Latest Technology Trends', 'Sharing my thoughts on the latest technology trends. What do you think?', NULL, '2024-02-22 09:20:00+00', 'public'),
('10000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000006', 'post', 'My First Post', 'My first post here. Looking forward to engaging with the community!', NULL, '2024-03-05 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000002', 'post', 'Best Practices for API Design', 'Community post: Let''s discuss best practices for API design.', NULL, '2024-03-10 14:15:00+00', 'community'),
('10000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000003', 'post', 'Project Progress Thoughts', 'Private thoughts on my current project progress.', NULL, '2024-03-12 16:30:00+00', 'private'),

-- Articles
('10000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000002', 'article', 'Introduction to Microservices Architecture', 'Introduction to Microservices Architecture

Microservices architecture has become increasingly popular in recent years. This article explores the key concepts, benefits, and challenges of building applications using microservices.

Key Benefits:
- Scalability
- Technology diversity
- Independent deployment
- Fault isolation

Challenges:
- Distributed system complexity
- Data consistency
- Network latency
- Service coordination', NULL, '2024-01-30 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000003', 'article', 'Best Practices for REST API Design', 'Best Practices for REST API Design

Designing a good REST API requires careful consideration of several factors. This article covers the most important principles.

1. Use proper HTTP methods
2. Follow RESTful conventions
3. Implement proper error handling
4. Use appropriate status codes
5. Version your API', NULL, '2024-02-12 11:30:00+00', 'public'),
('10000000-0000-0000-0000-000000000015', '00000000-0000-0000-0000-000000000001', 'article', 'Security Best Practices for Web Applications', 'Security Best Practices for Web Applications

Security should be a top priority in web development. This article discusses essential security practices.', NULL, '2024-02-25 13:45:00+00', 'community'),
('10000000-0000-0000-0000-000000000016', '00000000-0000-0000-0000-000000000004', 'article', 'Getting Started with Go Programming', 'Getting Started with Go Programming

Go is a powerful programming language developed by Google. This beginner-friendly guide covers the basics.', NULL, '2024-03-01 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000017', '00000000-0000-0000-0000-000000000002', 'article', 'Database Optimization Techniques', 'Database Optimization Techniques

Learn how to optimize your database queries and improve application performance.', NULL, '2024-03-08 15:20:00+00', 'public'),
('10000000-0000-0000-0000-000000000018', '00000000-0000-0000-0000-000000000005', 'article', 'My Personal Development Journey', 'My Personal Development Journey

A reflection on my growth as a developer over the past year.', NULL, '2024-03-14 10:10:00+00', 'private'),
('10000000-0000-0000-0000-000000000019', '00000000-0000-0000-0000-000000000003', 'article', 'Open Source Contributions', 'Community Article: Open Source Contributions

How contributing to open source projects can benefit your career.', NULL, '2024-03-15 12:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000020', '00000000-0000-0000-0000-000000000002', 'article', 'Advanced Topics in Distributed Systems', 'Advanced Topics in Distributed Systems

Exploring complex concepts in distributed computing.', NULL, '2024-03-18 14:30:00+00', 'public');

-- Additional PUBLICATIONS (80 publications to reach 100 total)
INSERT INTO publications (id, author_id, type, title, content, source, publication_date, visibility) VALUES
-- 21-40
('10000000-0000-0000-0000-000000000021', '00000000-0000-0000-0000-000000000002', 'post', 'Go Concurrency Notes', 'A few quick notes about goroutines, channels, and cancellation patterns.', NULL, '2024-03-19 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000022', '00000000-0000-0000-0000-000000000003', 'quote', 'Simplicity is the soul of efficiency.', 'Simplicity is the soul of efficiency.', 'Austin Freeman', '2024-03-19 10:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000023', '00000000-0000-0000-0000-000000000004', 'post', 'API Pagination Tips', 'Cursor-based pagination tends to scale better than offset-based in large datasets.', NULL, '2024-03-19 11:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000024', '00000000-0000-0000-0000-000000000005', 'article', 'Clean Architecture in Practice', 'A practical overview of layering, boundaries, and dependency inversion in backend services.', NULL, '2024-03-19 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000025', '00000000-0000-0000-0000-000000000006', 'quote', 'First make it work, then make it right, then make it fast.', 'First make it work, then make it right, then make it fast.', 'Kent Beck', '2024-03-19 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000026', '00000000-0000-0000-0000-000000000007', 'post', 'Reading List for Backend Engineers', 'Collected a short list of books and articles that helped me a lot.', NULL, '2024-03-19 14:00:00+00', 'private'),
('10000000-0000-0000-0000-000000000027', '00000000-0000-0000-0000-000000000002', 'article', 'Event-Driven Systems 101', 'When to use events, how to model them, and common pitfalls.', NULL, '2024-03-19 15:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000028', '00000000-0000-0000-0000-000000000003', 'post', 'Feature Flags and Rollouts', 'Gradual rollouts reduce risk; feature flags need ownership and cleanup.', NULL, '2024-03-19 16:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000029', '00000000-0000-0000-0000-000000000004', 'quote', 'Programs must be written for people to read.', 'Programs must be written for people to read.', 'Harold Abelson', '2024-03-19 17:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000030', '00000000-0000-0000-0000-000000000005', 'post', 'Indexing Strategy Basics', 'Indexes speed reads but slow writes; measure and keep them intentional.', NULL, '2024-03-19 18:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000031', '00000000-0000-0000-0000-000000000006', 'article', 'Intro to SQL Query Plans', 'How to read query plans and spot missing indexes and bad estimates.', NULL, '2024-03-20 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000032', '00000000-0000-0000-0000-000000000007', 'quote', 'What gets measured gets improved.', 'What gets measured gets improved.', 'Peter Drucker', '2024-03-20 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000033', '00000000-0000-0000-0000-000000000001', 'post', 'Admin Checklist for Production', 'Before deploying: backups, migrations, monitoring, and rollback steps.', NULL, '2024-03-20 11:00:00+00', 'private'),
('10000000-0000-0000-0000-000000000034', '00000000-0000-0000-0000-000000000002', 'quote', 'Make it correct, make it clear, make it concise.', 'Make it correct, make it clear, make it concise.', 'Anonymous', '2024-03-20 12:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000035', '00000000-0000-0000-0000-000000000003', 'article', 'Designing Stable APIs', 'Compatibility, versioning strategies, and deprecation policies that work.', NULL, '2024-03-20 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000036', '00000000-0000-0000-0000-000000000004', 'post', 'Testing HTTP Handlers', 'Table-driven tests plus golden responses keep handlers stable over time.', NULL, '2024-03-20 14:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000037', '00000000-0000-0000-0000-000000000005', 'quote', 'If you cannot explain it simply, you do not understand it well enough.', 'If you cannot explain it simply, you do not understand it well enough.', 'Albert Einstein', '2024-03-20 15:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000038', '00000000-0000-0000-0000-000000000006', 'post', 'Small Refactors Pay Off', 'Incremental refactors keep code healthy without huge risky rewrites.', NULL, '2024-03-20 16:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000039', '00000000-0000-0000-0000-000000000007', 'article', 'Security: Storing Passwords', 'Use bcrypt/argon2, unique salts, and never log credentials.', NULL, '2024-03-20 17:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000040', '00000000-0000-0000-0000-000000000001', 'quote', 'The best code is no code at all.', 'The best code is no code at all.', 'Jeff Atwood', '2024-03-20 18:00:00+00', 'public'),

-- 41-60
('10000000-0000-0000-0000-000000000041', '00000000-0000-0000-0000-000000000002', 'post', 'Caching: What to Cache', 'Cache the expensive reads, but make invalidation explicit and test it.', NULL, '2024-03-21 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000042', '00000000-0000-0000-0000-000000000003', 'quote', 'Premature optimization is the root of all evil.', 'Premature optimization is the root of all evil.', 'Donald Knuth', '2024-03-21 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000043', '00000000-0000-0000-0000-000000000004', 'post', 'Logging Levels Explained', 'Use debug for internals, info for lifecycle, warn for oddities, error for failures.', NULL, '2024-03-21 11:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000044', '00000000-0000-0000-0000-000000000005', 'article', 'Observability Fundamentals', 'Metrics, logs, and traces complement each other; start with SLIs and SLOs.', NULL, '2024-03-21 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000045', '00000000-0000-0000-0000-000000000006', 'quote', 'Debugging is twice as hard as writing the code.', 'Debugging is twice as hard as writing the code.', 'Brian Kernighan', '2024-03-21 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000046', '00000000-0000-0000-0000-000000000007', 'post', 'Learning PostgreSQL', 'Constraints are features. Start with a good schema, then optimize queries.', NULL, '2024-03-21 14:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000047', '00000000-0000-0000-0000-000000000001', 'article', 'Role-Based Access Control', 'RBAC needs clear policy rules and consistent enforcement at boundaries.', NULL, '2024-03-21 15:00:00+00', 'private'),
('10000000-0000-0000-0000-000000000048', '00000000-0000-0000-0000-000000000002', 'quote', 'Simplicity is prerequisite for reliability.', 'Simplicity is prerequisite for reliability.', 'Edsger Dijkstra', '2024-03-21 16:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000049', '00000000-0000-0000-0000-000000000003', 'post', 'Code Review Checklist', 'Naming, error handling, logging, tests, and performance are my core review points.', NULL, '2024-03-21 17:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000050', '00000000-0000-0000-0000-000000000004', 'article', 'Building Search Features', 'Start with simple text search, then add ranking, filters, and synonyms.', NULL, '2024-03-21 18:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000051', '00000000-0000-0000-0000-000000000005', 'post', 'Docker for Local Dev', 'Containers help keep envs consistent; document volumes and ports.', NULL, '2024-03-22 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000052', '00000000-0000-0000-0000-000000000006', 'quote', 'Quality is not an act, it is a habit.', 'Quality is not an act, it is a habit.', 'Aristotle', '2024-03-22 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000053', '00000000-0000-0000-0000-000000000007', 'post', 'HTTP Status Codes I Use Most', '200, 201, 204, 400, 401, 403, 404, 409, 422, 500 are enough for many APIs.', NULL, '2024-03-22 11:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000054', '00000000-0000-0000-0000-000000000001', 'quote', 'Any fool can write code that a computer can understand.', 'Any fool can write code that a computer can understand.', 'Martin Fowler', '2024-03-22 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000055', '00000000-0000-0000-0000-000000000002', 'article', 'Database Constraints as Tests', 'Constraints catch whole classes of bugs before they reach application code.', NULL, '2024-03-22 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000056', '00000000-0000-0000-0000-000000000003', 'post', 'Idempotency for APIs', 'Idempotent handlers simplify retries and reduce duplicated writes.', NULL, '2024-03-22 14:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000057', '00000000-0000-0000-0000-000000000004', 'quote', 'Talk is cheap. Show me the code.', 'Talk is cheap. Show me the code.', 'Linus Torvalds', '2024-03-22 15:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000058', '00000000-0000-0000-0000-000000000005', 'post', 'Repository Pattern Notes', 'Repositories should hide persistence details and keep business logic out.', NULL, '2024-03-22 16:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000059', '00000000-0000-0000-0000-000000000006', 'article', 'Graceful Shutdown in Go', 'Use contexts, timeouts, and close listeners first, then wait for in-flight work.', NULL, '2024-03-22 17:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000060', '00000000-0000-0000-0000-000000000007', 'quote', 'Do not repeat yourself.', 'Do not repeat yourself.', 'DRY Principle', '2024-03-22 18:00:00+00', 'public'),

-- 61-80
('10000000-0000-0000-0000-000000000061', '00000000-0000-0000-0000-000000000001', 'post', 'Incident Postmortems', 'Blameless postmortems focus on systems and improvements, not individuals.', NULL, '2024-03-23 09:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000062', '00000000-0000-0000-0000-000000000002', 'quote', 'When in doubt, use brute force.', 'When in doubt, use brute force.', 'Ken Thompson', '2024-03-23 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000063', '00000000-0000-0000-0000-000000000003', 'post', 'Schema Migrations', 'Keep migrations small, reversible when possible, and tested in CI.', NULL, '2024-03-23 11:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000064', '00000000-0000-0000-0000-000000000004', 'article', 'Designing for Failure', 'Time out, retry with backoff, and make partial failure a first-class scenario.', NULL, '2024-03-23 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000065', '00000000-0000-0000-0000-000000000005', 'quote', 'Make each program do one thing well.', 'Make each program do one thing well.', 'Unix Philosophy', '2024-03-23 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000066', '00000000-0000-0000-0000-000000000006', 'post', 'Tracing Requests', 'Propagate request IDs and context to make debugging distributed flows easier.', NULL, '2024-03-23 14:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000067', '00000000-0000-0000-0000-000000000007', 'article', 'Web Security Basics', 'CSRF, XSS, SSRF, and authz mistakes are the most common sources of incidents.', NULL, '2024-03-23 15:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000068', '00000000-0000-0000-0000-000000000001', 'quote', 'Build systems that make it hard to do the wrong thing.', 'Build systems that make it hard to do the wrong thing.', 'Anonymous', '2024-03-23 16:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000069', '00000000-0000-0000-0000-000000000002', 'post', 'Dependency Management', 'Pin versions, update regularly, and avoid unreviewed transitive changes.', NULL, '2024-03-23 17:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000070', '00000000-0000-0000-0000-000000000003', 'article', 'Text Search vs Full-Text Search', 'Know when LIKE is enough and when you need proper indexing and ranking.', NULL, '2024-03-23 18:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000071', '00000000-0000-0000-0000-000000000004', 'post', 'Designing JSON APIs', 'Use stable field names, clear errors, and avoid leaking internal DB schema.', NULL, '2024-03-24 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000072', '00000000-0000-0000-0000-000000000005', 'quote', 'Code never lies, comments sometimes do.', 'Code never lies, comments sometimes do.', 'Ron Jeffries', '2024-03-24 10:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000073', '00000000-0000-0000-0000-000000000006', 'post', 'Handling Time Zones', 'Store timestamps in UTC and convert at the edges; be consistent everywhere.', NULL, '2024-03-24 11:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000074', '00000000-0000-0000-0000-000000000007', 'article', 'Writing Maintainable SQL', 'Prefer explicit joins, keep queries readable, and measure with EXPLAIN.', NULL, '2024-03-24 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000075', '00000000-0000-0000-0000-000000000001', 'quote', 'Simple is better than complex.', 'Simple is better than complex.', 'Zen of Python', '2024-03-24 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000076', '00000000-0000-0000-0000-000000000002', 'post', 'Retries and Backoff', 'Retries should be bounded, jittered, and only for safe operations.', NULL, '2024-03-24 14:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000077', '00000000-0000-0000-0000-000000000003', 'quote', 'Perfect is the enemy of good.', 'Perfect is the enemy of good.', 'Voltaire', '2024-03-24 15:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000078', '00000000-0000-0000-0000-000000000004', 'post', 'Error Handling in Go', 'Wrap errors with context, avoid panics in servers, and keep messages actionable.', NULL, '2024-03-24 16:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000079', '00000000-0000-0000-0000-000000000005', 'article', 'Choosing Data Structures', 'Start from access patterns: maps for lookup, slices for order, and keep it simple.', NULL, '2024-03-24 17:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000080', '00000000-0000-0000-0000-000000000006', 'quote', 'The only real mistake is the one from which we learn nothing.', 'The only real mistake is the one from which we learn nothing.', 'Henry Ford', '2024-03-24 18:00:00+00', 'public'),

-- 81-100
('10000000-0000-0000-0000-000000000081', '00000000-0000-0000-0000-000000000007', 'post', 'Docs as Code', 'Version documentation with code so changes get reviewed and tested.', NULL, '2024-03-25 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000082', '00000000-0000-0000-0000-000000000001', 'quote', 'It always seems impossible until it is done.', 'It always seems impossible until it is done.', 'Nelson Mandela', '2024-03-25 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000083', '00000000-0000-0000-0000-000000000002', 'post', 'Background Jobs', 'Prefer queues for long work; keep request handlers fast and predictable.', NULL, '2024-03-25 11:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000084', '00000000-0000-0000-0000-000000000003', 'article', 'Practical Rate Limiting', 'Token buckets and sliding windows are common; pick what fits your traffic.', NULL, '2024-03-25 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000085', '00000000-0000-0000-0000-000000000004', 'quote', 'Without data you are just another person with an opinion.', 'Without data you are just another person with an opinion.', 'W. Edwards Deming', '2024-03-25 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000086', '00000000-0000-0000-0000-000000000005', 'post', 'When to Use Transactions', 'Use transactions for consistency boundaries, not as a general locking tool.', NULL, '2024-03-25 14:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000087', '00000000-0000-0000-0000-000000000006', 'article', 'Scaling Read Traffic', 'Caching, replicas, and pre-computation each help in different ways.', NULL, '2024-03-25 15:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000088', '00000000-0000-0000-0000-000000000007', 'quote', 'A year from now you may wish you had started today.', 'A year from now you may wish you had started today.', 'Karen Lamb', '2024-03-25 16:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000089', '00000000-0000-0000-0000-000000000001', 'post', 'Keeping Backwards Compatibility', 'Additive changes are safer than breaking changes; use deprecations thoughtfully.', NULL, '2024-03-25 17:00:00+00', 'private'),
('10000000-0000-0000-0000-000000000090', '00000000-0000-0000-0000-000000000002', 'quote', 'The details are not the details. They make the design.', 'The details are not the details. They make the design.', 'Charles Eames', '2024-03-25 18:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000091', '00000000-0000-0000-0000-000000000003', 'post', 'Structuring Go Packages', 'Keep packages small, cohesive, and named by domain concepts.', NULL, '2024-03-26 09:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000092', '00000000-0000-0000-0000-000000000004', 'article', 'Handling User-Generated Content', 'Validate, sanitize, and moderate. Store raw and derived forms when needed.', NULL, '2024-03-26 10:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000093', '00000000-0000-0000-0000-000000000005', 'quote', 'The most reliable code is the code you do not run.', 'The most reliable code is the code you do not run.', 'Anonymous', '2024-03-26 11:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000094', '00000000-0000-0000-0000-000000000006', 'post', 'Choosing HTTP Methods', 'Use GET for reads, POST for creates, PUT/PATCH for updates, DELETE for deletes.', NULL, '2024-03-26 12:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000095', '00000000-0000-0000-0000-000000000007', 'article', 'Designing Notification Systems', 'Use idempotency, deduplication, and user preferences from day one.', NULL, '2024-03-26 13:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000096', '00000000-0000-0000-0000-000000000001', 'quote', 'The best way to predict the future is to invent it.', 'The best way to predict the future is to invent it.', 'Alan Kay', '2024-03-26 14:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000097', '00000000-0000-0000-0000-000000000002', 'post', 'Small Wins', 'Shipping small improvements weekly beats waiting for a perfect big release.', NULL, '2024-03-26 15:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000098', '00000000-0000-0000-0000-000000000003', 'quote', 'The best error message is the one that never shows up.', 'The best error message is the one that never shows up.', 'Thomas Fuchs', '2024-03-26 16:00:00+00', 'public'),
('10000000-0000-0000-0000-000000000099', '00000000-0000-0000-0000-000000000004', 'post', 'Consistency in Naming', 'Consistent naming reduces cognitive load and makes APIs easier to use.', NULL, '2024-03-26 17:00:00+00', 'community'),
('10000000-0000-0000-0000-000000000100', '00000000-0000-0000-0000-000000000005', 'article', 'Building Reliable Backends', 'Timeouts, retries, limits, and good defaults are the foundation of reliability.', NULL, '2024-03-26 18:00:00+00', 'public');

-- COMMENTS (mix of top-level and nested comments)
INSERT INTO comments (id, publication_id, parent_id, author_id, text, created_at) VALUES
-- Comments on publication 1
('20000000-0000-0000-0000-000000000001', '10000000-0000-0000-0000-000000000001', NULL, '00000000-0000-0000-0000-000000000004', 'Great quote! Very inspiring.', '2024-01-20 11:00:00+00'),
('20000000-0000-0000-0000-000000000002', '10000000-0000-0000-0000-000000000001', NULL, '00000000-0000-0000-0000-000000000005', 'I completely agree with this sentiment.', '2024-01-20 12:30:00+00'),
('20000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000001', '20000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000006', 'Same here! This quote motivates me every day.', '2024-01-20 13:15:00+00'),

-- Comments on publication 6 (post)
('20000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000006', NULL, '00000000-0000-0000-0000-000000000003', 'Which book did you read? I''m always looking for recommendations.', '2024-02-05 14:00:00+00'),
('20000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000006', NULL, '00000000-0000-0000-0000-000000000004', 'Clean Code by Robert Martin is a classic!', '2024-02-05 15:30:00+00'),
('20000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000002', 'I recommend "Designing Data-Intensive Applications" as well.', '2024-02-05 16:00:00+00'),
('20000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000005', 'Thanks for the recommendation!', '2024-02-05 17:00:00+00'),

-- Comments on publication 13 (article)
('20000000-0000-0000-0000-000000000008', '10000000-0000-0000-0000-000000000013', NULL, '00000000-0000-0000-0000-000000000003', 'Excellent article! Microservices are indeed complex but worth it.', '2024-01-30 11:00:00+00'),
('20000000-0000-0000-0000-000000000009', '10000000-0000-0000-0000-000000000013', NULL, '00000000-0000-0000-0000-000000000004', 'I have a question about service coordination. Can you elaborate?', '2024-01-30 12:00:00+00'),
('20000000-0000-0000-0000-000000000010', '10000000-0000-0000-0000-000000000013', '20000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000002', 'Service coordination can be handled using message queues or event-driven architectures.', '2024-01-30 13:00:00+00'),
('20000000-0000-0000-0000-000000000011', '10000000-0000-0000-0000-000000000013', NULL, '00000000-0000-0000-0000-000000000005', 'Great overview of the topic!', '2024-01-30 14:00:00+00'),

-- Comments on publication 7
('20000000-0000-0000-0000-000000000012', '10000000-0000-0000-0000-000000000007', NULL, '00000000-0000-0000-0000-000000000004', 'Looking forward to seeing what you build!', '2024-02-10 11:00:00+00'),
('20000000-0000-0000-0000-000000000013', '10000000-0000-0000-0000-000000000007', NULL, '00000000-0000-0000-0000-000000000006', 'AI + web development sounds interesting!', '2024-02-10 12:00:00+00'),

-- Comments on publication 14
('20000000-0000-0000-0000-000000000014', '10000000-0000-0000-0000-000000000014', NULL, '00000000-0000-0000-0000-000000000002', 'Very practical advice. Thanks for sharing!', '2024-02-12 12:00:00+00'),
('20000000-0000-0000-0000-000000000015', '10000000-0000-0000-0000-000000000014', NULL, '00000000-0000-0000-0000-000000000005', 'I''ve been following these practices and they work great.', '2024-02-12 13:00:00+00');

-- Additional COMMENTS for publications 21-100 (moderate: 2 each, + extra replies on every 10th)
INSERT INTO comments (id, publication_id, parent_id, author_id, text, created_at) VALUES
-- Publications 21-30
('20000000-0000-0000-0000-000000000016', '10000000-0000-0000-0000-000000000021', NULL, '00000000-0000-0000-0000-000000000004', 'Nice write-up. The cancellation part is often overlooked.', '2024-03-19 09:30:00+00'),
('20000000-0000-0000-0000-000000000017', '10000000-0000-0000-0000-000000000021', NULL, '00000000-0000-0000-0000-000000000006', 'Agree. Context propagation makes debugging much easier.', '2024-03-19 09:45:00+00'),
('20000000-0000-0000-0000-000000000018', '10000000-0000-0000-0000-000000000022', NULL, '00000000-0000-0000-0000-000000000005', 'Great quote, very applicable to engineering.', '2024-03-19 10:20:00+00'),
('20000000-0000-0000-0000-000000000019', '10000000-0000-0000-0000-000000000022', NULL, '00000000-0000-0000-0000-000000000007', 'Simplicity is underrated. Thanks for sharing.', '2024-03-19 10:40:00+00'),
('20000000-0000-0000-0000-000000000020', '10000000-0000-0000-0000-000000000023', NULL, '00000000-0000-0000-0000-000000000006', 'Cursor-based pagination saved us a lot of pain.', '2024-03-19 11:20:00+00'),
('20000000-0000-0000-0000-000000000021', '10000000-0000-0000-0000-000000000023', NULL, '00000000-0000-0000-0000-000000000004', 'Do you also include total counts or avoid them?', '2024-03-19 11:40:00+00'),
('20000000-0000-0000-0000-000000000022', '10000000-0000-0000-0000-000000000024', NULL, '00000000-0000-0000-0000-000000000003', 'I like the boundary-driven approach in clean architecture.', '2024-03-19 12:20:00+00'),
('20000000-0000-0000-0000-000000000023', '10000000-0000-0000-0000-000000000024', NULL, '00000000-0000-0000-0000-000000000005', 'This helped me rethink dependencies in my project.', '2024-03-19 12:40:00+00'),
('20000000-0000-0000-0000-000000000024', '10000000-0000-0000-0000-000000000025', NULL, '00000000-0000-0000-0000-000000000007', 'Classic advice. The order really matters.', '2024-03-19 13:20:00+00'),
('20000000-0000-0000-0000-000000000025', '10000000-0000-0000-0000-000000000025', NULL, '00000000-0000-0000-0000-000000000002', 'Also: measure before optimizing.', '2024-03-19 13:40:00+00'),
('20000000-0000-0000-0000-000000000026', '10000000-0000-0000-0000-000000000026', NULL, '00000000-0000-0000-0000-000000000004', 'Interested in that list. Any top 3 picks?', '2024-03-19 14:20:00+00'),
('20000000-0000-0000-0000-000000000027', '10000000-0000-0000-0000-000000000026', NULL, '00000000-0000-0000-0000-000000000006', 'I would add a book on databases as well.', '2024-03-19 14:35:00+00'),
('20000000-0000-0000-0000-000000000028', '10000000-0000-0000-0000-000000000027', NULL, '00000000-0000-0000-0000-000000000005', 'Events make integrations smoother when done carefully.', '2024-03-19 15:20:00+00'),
('20000000-0000-0000-0000-000000000029', '10000000-0000-0000-0000-000000000027', NULL, '00000000-0000-0000-0000-000000000003', 'How do you handle ordering guarantees?', '2024-03-19 15:40:00+00'),
('20000000-0000-0000-0000-000000000030', '10000000-0000-0000-0000-000000000028', NULL, '00000000-0000-0000-0000-000000000004', 'Feature flags are great, but cleanup is important.', '2024-03-19 16:20:00+00'),
('20000000-0000-0000-0000-000000000031', '10000000-0000-0000-0000-000000000028', NULL, '00000000-0000-0000-0000-000000000007', 'Agree. Old flags become technical debt quickly.', '2024-03-19 16:40:00+00'),
('20000000-0000-0000-0000-000000000032', '10000000-0000-0000-0000-000000000029', NULL, '00000000-0000-0000-0000-000000000006', 'This is one of my favorite quotes about code.', '2024-03-19 17:20:00+00'),
('20000000-0000-0000-0000-000000000033', '10000000-0000-0000-0000-000000000029', NULL, '00000000-0000-0000-0000-000000000005', 'Readable code pays off in the long run.', '2024-03-19 17:40:00+00'),
('20000000-0000-0000-0000-000000000034', '10000000-0000-0000-0000-000000000030', NULL, '00000000-0000-0000-0000-000000000003', 'Index bloat is real. Pruning unused indexes helps.', '2024-03-19 18:20:00+00'),
('20000000-0000-0000-0000-000000000035', '10000000-0000-0000-0000-000000000030', NULL, '00000000-0000-0000-0000-000000000004', 'Totally. Measure write amplification too.', '2024-03-19 18:40:00+00'),

-- Publications 31-40
('20000000-0000-0000-0000-000000000036', '10000000-0000-0000-0000-000000000031', NULL, '00000000-0000-0000-0000-000000000004', 'Learning to read plans is a huge skill boost.', '2024-03-20 09:20:00+00'),
('20000000-0000-0000-0000-000000000037', '10000000-0000-0000-0000-000000000031', NULL, '00000000-0000-0000-0000-000000000002', 'EXPLAIN ANALYZE is my favorite tool.', '2024-03-20 09:40:00+00'),
('20000000-0000-0000-0000-000000000038', '10000000-0000-0000-0000-000000000032', NULL, '00000000-0000-0000-0000-000000000005', 'Metrics really drive better decisions.', '2024-03-20 10:20:00+00'),
('20000000-0000-0000-0000-000000000039', '10000000-0000-0000-0000-000000000032', NULL, '00000000-0000-0000-0000-000000000006', 'Agree. Observability first, optimizations second.', '2024-03-20 10:40:00+00'),
('20000000-0000-0000-0000-000000000040', '10000000-0000-0000-0000-000000000033', NULL, '00000000-0000-0000-0000-000000000007', 'Good checklist. Rollback plans are often missing.', '2024-03-20 11:20:00+00'),
('20000000-0000-0000-0000-000000000041', '10000000-0000-0000-0000-000000000033', NULL, '00000000-0000-0000-0000-000000000001', 'Exactly. Preparedness beats heroics.', '2024-03-20 11:45:00+00'),
('20000000-0000-0000-0000-000000000042', '10000000-0000-0000-0000-000000000034', NULL, '00000000-0000-0000-0000-000000000003', 'Clarity is a feature.', '2024-03-20 12:20:00+00'),
('20000000-0000-0000-0000-000000000043', '10000000-0000-0000-0000-000000000034', NULL, '00000000-0000-0000-0000-000000000004', 'Concise code is easier to maintain.', '2024-03-20 12:40:00+00'),
('20000000-0000-0000-0000-000000000044', '10000000-0000-0000-0000-000000000035', NULL, '00000000-0000-0000-0000-000000000005', 'Deprecation policies save customers from surprises.', '2024-03-20 13:20:00+00'),
('20000000-0000-0000-0000-000000000045', '10000000-0000-0000-0000-000000000035', NULL, '00000000-0000-0000-0000-000000000006', 'Versioning is hard, but necessary.', '2024-03-20 13:40:00+00'),
('20000000-0000-0000-0000-000000000046', '10000000-0000-0000-0000-000000000036', NULL, '00000000-0000-0000-0000-000000000002', 'Golden tests are great for response stability.', '2024-03-20 14:20:00+00'),
('20000000-0000-0000-0000-000000000047', '10000000-0000-0000-0000-000000000036', NULL, '00000000-0000-0000-0000-000000000004', 'Table-driven tests keep things clean.', '2024-03-20 14:40:00+00'),
('20000000-0000-0000-0000-000000000048', '10000000-0000-0000-0000-000000000037', NULL, '00000000-0000-0000-0000-000000000007', 'Explaining clearly is a real skill.', '2024-03-20 15:20:00+00'),
('20000000-0000-0000-0000-000000000049', '10000000-0000-0000-0000-000000000037', NULL, '00000000-0000-0000-0000-000000000003', 'Simple explanations usually reveal gaps in thinking.', '2024-03-20 15:40:00+00'),
('20000000-0000-0000-0000-000000000050', '10000000-0000-0000-0000-000000000038', NULL, '00000000-0000-0000-0000-000000000005', 'Incremental refactors are the best kind.', '2024-03-20 16:20:00+00'),
('20000000-0000-0000-0000-000000000051', '10000000-0000-0000-0000-000000000038', NULL, '00000000-0000-0000-0000-000000000006', 'Agree. Big bang rewrites are risky.', '2024-03-20 16:40:00+00'),
('20000000-0000-0000-0000-000000000052', '10000000-0000-0000-0000-000000000039', NULL, '00000000-0000-0000-0000-000000000004', 'Good reminder: never log secrets.', '2024-03-20 17:20:00+00'),
('20000000-0000-0000-0000-000000000053', '10000000-0000-0000-0000-000000000039', NULL, '00000000-0000-0000-0000-000000000002', 'And always rate-limit login endpoints.', '2024-03-20 17:40:00+00'),
('20000000-0000-0000-0000-000000000054', '10000000-0000-0000-0000-000000000040', NULL, '00000000-0000-0000-0000-000000000005', 'Less code, fewer bugs.', '2024-03-20 18:20:00+00'),
('20000000-0000-0000-0000-000000000055', '10000000-0000-0000-0000-000000000040', NULL, '00000000-0000-0000-0000-000000000007', 'Sometimes the best feature is removing features.', '2024-03-20 18:40:00+00'),

-- Publications 41-50 (extra reply for publication 50)
('20000000-0000-0000-0000-000000000056', '10000000-0000-0000-0000-000000000041', NULL, '00000000-0000-0000-0000-000000000004', 'Caching is powerful, but invalidation is tricky.', '2024-03-21 09:20:00+00'),
('20000000-0000-0000-0000-000000000057', '10000000-0000-0000-0000-000000000041', NULL, '00000000-0000-0000-0000-000000000006', 'Do you prefer write-through or cache-aside?', '2024-03-21 09:40:00+00'),
('20000000-0000-0000-0000-000000000058', '10000000-0000-0000-0000-000000000042', NULL, '00000000-0000-0000-0000-000000000007', 'Still relevant advice today.', '2024-03-21 10:20:00+00'),
('20000000-0000-0000-0000-000000000059', '10000000-0000-0000-0000-000000000042', NULL, '00000000-0000-0000-0000-000000000003', 'Agree. Optimize after you understand the bottleneck.', '2024-03-21 10:40:00+00'),
('20000000-0000-0000-0000-000000000060', '10000000-0000-0000-0000-000000000043', NULL, '00000000-0000-0000-0000-000000000005', 'Levels help keep logs useful.', '2024-03-21 11:20:00+00'),
('20000000-0000-0000-0000-000000000061', '10000000-0000-0000-0000-000000000043', NULL, '00000000-0000-0000-0000-000000000004', 'And sampling helps at high volume.', '2024-03-21 11:40:00+00'),
('20000000-0000-0000-0000-000000000062', '10000000-0000-0000-0000-000000000044', NULL, '00000000-0000-0000-0000-000000000006', 'Starting from SLOs is the right move.', '2024-03-21 12:20:00+00'),
('20000000-0000-0000-0000-000000000063', '10000000-0000-0000-0000-000000000044', NULL, '00000000-0000-0000-0000-000000000002', 'Traces reveal where latency really comes from.', '2024-03-21 12:40:00+00'),
('20000000-0000-0000-0000-000000000064', '10000000-0000-0000-0000-000000000045', NULL, '00000000-0000-0000-0000-000000000003', 'So true. The first version is never perfect.', '2024-03-21 13:20:00+00'),
('20000000-0000-0000-0000-000000000065', '10000000-0000-0000-0000-000000000045', NULL, '00000000-0000-0000-0000-000000000005', 'And tests reduce debugging time.', '2024-03-21 13:40:00+00'),
('20000000-0000-0000-0000-000000000066', '10000000-0000-0000-0000-000000000046', NULL, '00000000-0000-0000-0000-000000000004', 'Constraints prevent many bugs by design.', '2024-03-21 14:20:00+00'),
('20000000-0000-0000-0000-000000000067', '10000000-0000-0000-0000-000000000046', NULL, '00000000-0000-0000-0000-000000000006', 'PostgreSQL has great features for this.', '2024-03-21 14:40:00+00'),
('20000000-0000-0000-0000-000000000068', '10000000-0000-0000-0000-000000000047', NULL, '00000000-0000-0000-0000-000000000002', 'Policy clarity makes enforcement easier.', '2024-03-21 15:20:00+00'),
('20000000-0000-0000-0000-000000000069', '10000000-0000-0000-0000-000000000047', NULL, '00000000-0000-0000-0000-000000000001', 'Agree. Keep checks close to boundaries.', '2024-03-21 15:40:00+00'),
('20000000-0000-0000-0000-000000000070', '10000000-0000-0000-0000-000000000048', NULL, '00000000-0000-0000-0000-000000000005', 'Reliability is easier when systems are simple.', '2024-03-21 16:20:00+00'),
('20000000-0000-0000-0000-000000000071', '10000000-0000-0000-0000-000000000048', NULL, '00000000-0000-0000-0000-000000000007', 'Less moving parts, fewer outages.', '2024-03-21 16:40:00+00'),
('20000000-0000-0000-0000-000000000072', '10000000-0000-0000-0000-000000000049', NULL, '00000000-0000-0000-0000-000000000004', 'Checklist reviews help keep teams aligned.', '2024-03-21 17:20:00+00'),
('20000000-0000-0000-0000-000000000073', '10000000-0000-0000-0000-000000000049', NULL, '00000000-0000-0000-0000-000000000006', 'Agree. And they teach juniors a lot.', '2024-03-21 17:40:00+00'),
('20000000-0000-0000-0000-000000000074', '10000000-0000-0000-0000-000000000050', NULL, '00000000-0000-0000-0000-000000000003', 'Search is deceptively complex once you add ranking.', '2024-03-21 18:20:00+00'),
('20000000-0000-0000-0000-000000000075', '10000000-0000-0000-0000-000000000050', '20000000-0000-0000-0000-000000000074', '00000000-0000-0000-0000-000000000002', 'True. Start simple and iterate with real queries.', '2024-03-21 18:35:00+00'),
('20000000-0000-0000-0000-000000000076', '10000000-0000-0000-0000-000000000050', NULL, '00000000-0000-0000-0000-000000000005', 'Filters and facets are what users ask for first.', '2024-03-21 18:45:00+00'),

-- Publications 51-60 (extra reply for publication 60)
('20000000-0000-0000-0000-000000000077', '10000000-0000-0000-0000-000000000051', NULL, '00000000-0000-0000-0000-000000000006', 'Dockerfiles are documentation too.', '2024-03-22 09:20:00+00'),
('20000000-0000-0000-0000-000000000078', '10000000-0000-0000-0000-000000000051', NULL, '00000000-0000-0000-0000-000000000004', 'Volume usage and port mapping always matter.', '2024-03-22 09:40:00+00'),
('20000000-0000-0000-0000-000000000079', '10000000-0000-0000-0000-000000000052', NULL, '00000000-0000-0000-0000-000000000005', 'Consistency builds great products.', '2024-03-22 10:20:00+00'),
('20000000-0000-0000-0000-000000000080', '10000000-0000-0000-0000-000000000052', NULL, '00000000-0000-0000-0000-000000000007', 'Habits compound over time.', '2024-03-22 10:40:00+00'),
('20000000-0000-0000-0000-000000000081', '10000000-0000-0000-0000-000000000053', NULL, '00000000-0000-0000-0000-000000000003', '422 is underrated for validation errors.', '2024-03-22 11:20:00+00'),
('20000000-0000-0000-0000-000000000082', '10000000-0000-0000-0000-000000000053', NULL, '00000000-0000-0000-0000-000000000006', 'And 409 is great for conflicts.', '2024-03-22 11:40:00+00'),
('20000000-0000-0000-0000-000000000083', '10000000-0000-0000-0000-000000000054', NULL, '00000000-0000-0000-0000-000000000004', 'Readable code helps teams scale.', '2024-03-22 12:20:00+00'),
('20000000-0000-0000-0000-000000000084', '10000000-0000-0000-0000-000000000054', NULL, '00000000-0000-0000-0000-000000000005', 'Agree. Humans are the real runtime.', '2024-03-22 12:40:00+00'),
('20000000-0000-0000-0000-000000000085', '10000000-0000-0000-0000-000000000055', NULL, '00000000-0000-0000-0000-000000000002', 'Constraints are great guardrails.', '2024-03-22 13:20:00+00'),
('20000000-0000-0000-0000-000000000086', '10000000-0000-0000-0000-000000000055', NULL, '00000000-0000-0000-0000-000000000006', 'Also helps keep data clean for analytics.', '2024-03-22 13:40:00+00'),
('20000000-0000-0000-0000-000000000087', '10000000-0000-0000-0000-000000000056', NULL, '00000000-0000-0000-0000-000000000004', 'Idempotency keys are a life saver.', '2024-03-22 14:20:00+00'),
('20000000-0000-0000-0000-000000000088', '10000000-0000-0000-0000-000000000056', NULL, '00000000-0000-0000-0000-000000000007', 'Retries become much safer with idempotency.', '2024-03-22 14:40:00+00'),
('20000000-0000-0000-0000-000000000089', '10000000-0000-0000-0000-000000000057', NULL, '00000000-0000-0000-0000-000000000005', 'The quote is timeless.', '2024-03-22 15:20:00+00'),
('20000000-0000-0000-0000-000000000090', '10000000-0000-0000-0000-000000000057', NULL, '00000000-0000-0000-0000-000000000003', 'Showing code clarifies intent.', '2024-03-22 15:40:00+00'),
('20000000-0000-0000-0000-000000000091', '10000000-0000-0000-0000-000000000058', NULL, '00000000-0000-0000-0000-000000000006', 'Repositories keep domain logic clean.', '2024-03-22 16:20:00+00'),
('20000000-0000-0000-0000-000000000092', '10000000-0000-0000-0000-000000000058', NULL, '00000000-0000-0000-0000-000000000004', 'Agree. Hide SQL details behind an interface.', '2024-03-22 16:40:00+00'),
('20000000-0000-0000-0000-000000000093', '10000000-0000-0000-0000-000000000059', NULL, '00000000-0000-0000-0000-000000000002', 'Shutdown issues are hard to reproduce without tests.', '2024-03-22 17:20:00+00'),
('20000000-0000-0000-0000-000000000094', '10000000-0000-0000-0000-000000000059', NULL, '00000000-0000-0000-0000-000000000006', 'Time-bounded shutdown is key for deployments.', '2024-03-22 17:40:00+00'),
('20000000-0000-0000-0000-000000000095', '10000000-0000-0000-0000-000000000060', NULL, '00000000-0000-0000-0000-000000000005', 'DRY is great when it reduces repetition, not clarity.', '2024-03-22 18:20:00+00'),
('20000000-0000-0000-0000-000000000096', '10000000-0000-0000-0000-000000000060', '20000000-0000-0000-0000-000000000095', '00000000-0000-0000-0000-000000000003', 'Exactly. Some duplication is better than the wrong abstraction.', '2024-03-22 18:35:00+00'),
('20000000-0000-0000-0000-000000000097', '10000000-0000-0000-0000-000000000060', NULL, '00000000-0000-0000-0000-000000000007', 'Balance is everything.', '2024-03-22 18:45:00+00'),

-- Publications 61-70 (extra reply for publication 70)
('20000000-0000-0000-0000-000000000098', '10000000-0000-0000-0000-000000000061', NULL, '00000000-0000-0000-0000-000000000004', 'Blameless is the only sustainable approach.', '2024-03-23 09:20:00+00'),
('20000000-0000-0000-0000-000000000099', '10000000-0000-0000-0000-000000000061', NULL, '00000000-0000-0000-0000-000000000006', 'Postmortems should produce action items with owners.', '2024-03-23 09:40:00+00'),
('20000000-0000-0000-0000-000000000100', '10000000-0000-0000-0000-000000000062', NULL, '00000000-0000-0000-0000-000000000005', 'Sometimes brute force is the simplest route.', '2024-03-23 10:20:00+00'),
('20000000-0000-0000-0000-000000000101', '10000000-0000-0000-0000-000000000062', NULL, '00000000-0000-0000-0000-000000000007', 'But keep it readable.', '2024-03-23 10:40:00+00'),
('20000000-0000-0000-0000-000000000102', '10000000-0000-0000-0000-000000000063', NULL, '00000000-0000-0000-0000-000000000003', 'Testing migrations in CI catches surprises.', '2024-03-23 11:20:00+00'),
('20000000-0000-0000-0000-000000000103', '10000000-0000-0000-0000-000000000063', NULL, '00000000-0000-0000-0000-000000000004', 'Small steps make rollbacks easier.', '2024-03-23 11:40:00+00'),
('20000000-0000-0000-0000-000000000104', '10000000-0000-0000-0000-000000000064', NULL, '00000000-0000-0000-0000-000000000002', 'Timeouts and backoff prevent cascading failures.', '2024-03-23 12:20:00+00'),
('20000000-0000-0000-0000-000000000105', '10000000-0000-0000-0000-000000000064', NULL, '00000000-0000-0000-0000-000000000006', 'And bulkheads help a lot.', '2024-03-23 12:40:00+00'),
('20000000-0000-0000-0000-000000000106', '10000000-0000-0000-0000-000000000065', NULL, '00000000-0000-0000-0000-000000000005', 'Single responsibility keeps tools sharp.', '2024-03-23 13:20:00+00'),
('20000000-0000-0000-0000-000000000107', '10000000-0000-0000-0000-000000000065', NULL, '00000000-0000-0000-0000-000000000007', 'Pipelines help combine them later.', '2024-03-23 13:40:00+00'),
('20000000-0000-0000-0000-000000000108', '10000000-0000-0000-0000-000000000066', NULL, '00000000-0000-0000-0000-000000000004', 'Request IDs are extremely helpful in logs.', '2024-03-23 14:20:00+00'),
('20000000-0000-0000-0000-000000000109', '10000000-0000-0000-0000-000000000066', NULL, '00000000-0000-0000-0000-000000000006', 'Traces + logs together are best.', '2024-03-23 14:40:00+00'),
('20000000-0000-0000-0000-000000000110', '10000000-0000-0000-0000-000000000067', NULL, '00000000-0000-0000-0000-000000000003', 'SSRF is still surprisingly common.', '2024-03-23 15:20:00+00'),
('20000000-0000-0000-0000-000000000111', '10000000-0000-0000-0000-000000000067', NULL, '00000000-0000-0000-0000-000000000002', 'Defense in depth is key.', '2024-03-23 15:40:00+00'),
('20000000-0000-0000-0000-000000000112', '10000000-0000-0000-0000-000000000068', NULL, '00000000-0000-0000-0000-000000000007', 'Guardrails help new contributors.', '2024-03-23 16:20:00+00'),
('20000000-0000-0000-0000-000000000113', '10000000-0000-0000-0000-000000000068', NULL, '00000000-0000-0000-0000-000000000005', 'Agreed. Make safe defaults the default.', '2024-03-23 16:40:00+00'),
('20000000-0000-0000-0000-000000000114', '10000000-0000-0000-0000-000000000069', NULL, '00000000-0000-0000-0000-000000000004', 'Pinning versions avoids surprise breakages.', '2024-03-23 17:20:00+00'),
('20000000-0000-0000-0000-000000000115', '10000000-0000-0000-0000-000000000069', NULL, '00000000-0000-0000-0000-000000000006', 'And dependabot-style updates keep it manageable.', '2024-03-23 17:40:00+00'),
('20000000-0000-0000-0000-000000000116', '10000000-0000-0000-0000-000000000070', NULL, '00000000-0000-0000-0000-000000000003', 'Full-text search is worth it once you need ranking.', '2024-03-23 18:20:00+00'),
('20000000-0000-0000-0000-000000000117', '10000000-0000-0000-0000-000000000070', '20000000-0000-0000-0000-000000000116', '00000000-0000-0000-0000-000000000005', 'Agreed. LIKE is fine for prototypes though.', '2024-03-23 18:35:00+00'),
('20000000-0000-0000-0000-000000000118', '10000000-0000-0000-0000-000000000070', NULL, '00000000-0000-0000-0000-000000000007', 'Indexing strategy makes or breaks search.', '2024-03-23 18:45:00+00'),

-- Publications 71-80 (extra reply for publication 80)
('20000000-0000-0000-0000-000000000119', '10000000-0000-0000-0000-000000000071', NULL, '00000000-0000-0000-0000-000000000006', 'Clear errors are the best UX for developers.', '2024-03-24 09:20:00+00'),
('20000000-0000-0000-0000-000000000120', '10000000-0000-0000-0000-000000000071', NULL, '00000000-0000-0000-0000-000000000004', 'And consistent schemas help client code.', '2024-03-24 09:40:00+00'),
('20000000-0000-0000-0000-000000000121', '10000000-0000-0000-0000-000000000072', NULL, '00000000-0000-0000-0000-000000000005', 'Comments should be maintained like code.', '2024-03-24 10:20:00+00'),
('20000000-0000-0000-0000-000000000122', '10000000-0000-0000-0000-000000000072', NULL, '00000000-0000-0000-0000-000000000007', 'Exactly. Keep them accurate or remove.', '2024-03-24 10:40:00+00'),
('20000000-0000-0000-0000-000000000123', '10000000-0000-0000-0000-000000000073', NULL, '00000000-0000-0000-0000-000000000003', 'UTC everywhere avoids many surprises.', '2024-03-24 11:20:00+00'),
('20000000-0000-0000-0000-000000000124', '10000000-0000-0000-0000-000000000073', NULL, '00000000-0000-0000-0000-000000000006', 'Daylight saving issues are painful.', '2024-03-24 11:40:00+00'),
('20000000-0000-0000-0000-000000000125', '10000000-0000-0000-0000-000000000074', NULL, '00000000-0000-0000-0000-000000000004', 'Readable SQL is a gift to future you.', '2024-03-24 12:20:00+00'),
('20000000-0000-0000-0000-000000000126', '10000000-0000-0000-0000-000000000074', NULL, '00000000-0000-0000-0000-000000000002', 'Agreed. Explicit joins are best.', '2024-03-24 12:40:00+00'),
('20000000-0000-0000-0000-000000000127', '10000000-0000-0000-0000-000000000075', NULL, '00000000-0000-0000-0000-000000000005', 'Simple designs age better.', '2024-03-24 13:20:00+00'),
('20000000-0000-0000-0000-000000000128', '10000000-0000-0000-0000-000000000075', NULL, '00000000-0000-0000-0000-000000000007', 'Complexity has a hidden cost.', '2024-03-24 13:40:00+00'),
('20000000-0000-0000-0000-000000000129', '10000000-0000-0000-0000-000000000076', NULL, '00000000-0000-0000-0000-000000000006', 'Jitter avoids synchronized retries.', '2024-03-24 14:20:00+00'),
('20000000-0000-0000-0000-000000000130', '10000000-0000-0000-0000-000000000076', NULL, '00000000-0000-0000-0000-000000000004', 'And always cap retry attempts.', '2024-03-24 14:40:00+00'),
('20000000-0000-0000-0000-000000000131', '10000000-0000-0000-0000-000000000077', NULL, '00000000-0000-0000-0000-000000000005', 'Very true. Good enough can ship.', '2024-03-24 15:20:00+00'),
('20000000-0000-0000-0000-000000000132', '10000000-0000-0000-0000-000000000077', NULL, '00000000-0000-0000-0000-000000000003', 'Perfection delays learning.', '2024-03-24 15:40:00+00'),
('20000000-0000-0000-0000-000000000133', '10000000-0000-0000-0000-000000000078', NULL, '00000000-0000-0000-0000-000000000002', 'Wrapping errors with context helps a lot.', '2024-03-24 16:20:00+00'),
('20000000-0000-0000-0000-000000000134', '10000000-0000-0000-0000-000000000078', NULL, '00000000-0000-0000-0000-000000000006', 'And centralize error-to-HTTP mapping.', '2024-03-24 16:40:00+00'),
('20000000-0000-0000-0000-000000000135', '10000000-0000-0000-0000-000000000079', NULL, '00000000-0000-0000-0000-000000000004', 'Access patterns should drive data structures.', '2024-03-24 17:20:00+00'),
('20000000-0000-0000-0000-000000000136', '10000000-0000-0000-0000-000000000079', NULL, '00000000-0000-0000-0000-000000000005', 'Start with simple choices and measure.', '2024-03-24 17:40:00+00'),
('20000000-0000-0000-0000-000000000137', '10000000-0000-0000-0000-000000000080', NULL, '00000000-0000-0000-0000-000000000007', 'Learning is the only way forward.', '2024-03-24 18:20:00+00'),
('20000000-0000-0000-0000-000000000138', '10000000-0000-0000-0000-000000000080', '20000000-0000-0000-0000-000000000137', '00000000-0000-0000-0000-000000000006', 'Agreed. Mistakes are feedback.', '2024-03-24 18:35:00+00'),
('20000000-0000-0000-0000-000000000139', '10000000-0000-0000-0000-000000000080', NULL, '00000000-0000-0000-0000-000000000003', 'A good team learns quickly.', '2024-03-24 18:45:00+00'),

-- Publications 81-90 (extra reply for publication 90)
('20000000-0000-0000-0000-000000000140', '10000000-0000-0000-0000-000000000081', NULL, '00000000-0000-0000-0000-000000000004', 'Docs in git make changes visible and reviewable.', '2024-03-25 09:20:00+00'),
('20000000-0000-0000-0000-000000000141', '10000000-0000-0000-0000-000000000081', NULL, '00000000-0000-0000-0000-000000000006', 'Also helps onboard new people faster.', '2024-03-25 09:40:00+00'),
('20000000-0000-0000-0000-000000000142', '10000000-0000-0000-0000-000000000082', NULL, '00000000-0000-0000-0000-000000000005', 'A motivating quote.', '2024-03-25 10:20:00+00'),
('20000000-0000-0000-0000-000000000143', '10000000-0000-0000-0000-000000000082', NULL, '00000000-0000-0000-0000-000000000007', 'Starting is the hardest part.', '2024-03-25 10:40:00+00'),
('20000000-0000-0000-0000-000000000144', '10000000-0000-0000-0000-000000000083', NULL, '00000000-0000-0000-0000-000000000003', 'Queues simplify load spikes dramatically.', '2024-03-25 11:20:00+00'),
('20000000-0000-0000-0000-000000000145', '10000000-0000-0000-0000-000000000083', NULL, '00000000-0000-0000-0000-000000000002', 'And they help with retries and scheduling.', '2024-03-25 11:40:00+00'),
('20000000-0000-0000-0000-000000000146', '10000000-0000-0000-0000-000000000084', NULL, '00000000-0000-0000-0000-000000000006', 'Rate limiting prevents abuse and accidents.', '2024-03-25 12:20:00+00'),
('20000000-0000-0000-0000-000000000147', '10000000-0000-0000-0000-000000000084', NULL, '00000000-0000-0000-0000-000000000004', 'Token bucket is a good default.', '2024-03-25 12:40:00+00'),
('20000000-0000-0000-0000-000000000148', '10000000-0000-0000-0000-000000000085', NULL, '00000000-0000-0000-0000-000000000005', 'Data beats opinions every time.', '2024-03-25 13:20:00+00'),
('20000000-0000-0000-0000-000000000149', '10000000-0000-0000-0000-000000000085', NULL, '00000000-0000-0000-0000-000000000007', 'Measure, then decide.', '2024-03-25 13:40:00+00'),
('20000000-0000-0000-0000-000000000150', '10000000-0000-0000-0000-000000000086', NULL, '00000000-0000-0000-0000-000000000004', 'Transactions define consistency boundaries.', '2024-03-25 14:20:00+00'),
('20000000-0000-0000-0000-000000000151', '10000000-0000-0000-0000-000000000086', NULL, '00000000-0000-0000-0000-000000000006', 'And they should stay short.', '2024-03-25 14:40:00+00'),
('20000000-0000-0000-0000-000000000152', '10000000-0000-0000-0000-000000000087', NULL, '00000000-0000-0000-0000-000000000003', 'Replicas help, but keep replication lag in mind.', '2024-03-25 15:20:00+00'),
('20000000-0000-0000-0000-000000000153', '10000000-0000-0000-0000-000000000087', NULL, '00000000-0000-0000-0000-000000000005', 'Pre-computation is underrated for hot endpoints.', '2024-03-25 15:40:00+00'),
('20000000-0000-0000-0000-000000000154', '10000000-0000-0000-0000-000000000088', NULL, '00000000-0000-0000-0000-000000000007', 'Starting today matters more than planning forever.', '2024-03-25 16:20:00+00'),
('20000000-0000-0000-0000-000000000155', '10000000-0000-0000-0000-000000000088', NULL, '00000000-0000-0000-0000-000000000006', 'Progress compounds.', '2024-03-25 16:40:00+00'),
('20000000-0000-0000-0000-000000000156', '10000000-0000-0000-0000-000000000089', NULL, '00000000-0000-0000-0000-000000000002', 'Compatibility is a promise to your users.', '2024-03-25 17:20:00+00'),
('20000000-0000-0000-0000-000000000157', '10000000-0000-0000-0000-000000000089', NULL, '00000000-0000-0000-0000-000000000001', 'Deprecations should be communicated early.', '2024-03-25 17:40:00+00'),
('20000000-0000-0000-0000-000000000158', '10000000-0000-0000-0000-000000000090', NULL, '00000000-0000-0000-0000-000000000004', 'Design is in the details indeed.', '2024-03-25 18:20:00+00'),
('20000000-0000-0000-0000-000000000159', '10000000-0000-0000-0000-000000000090', '20000000-0000-0000-0000-000000000158', '00000000-0000-0000-0000-000000000003', 'Small inconsistencies add up fast.', '2024-03-25 18:35:00+00'),
('20000000-0000-0000-0000-000000000160', '10000000-0000-0000-0000-000000000090', NULL, '00000000-0000-0000-0000-000000000005', 'This is a great reminder for API design too.', '2024-03-25 18:45:00+00'),

-- Publications 91-100 (extra reply for publication 100)
('20000000-0000-0000-0000-000000000161', '10000000-0000-0000-0000-000000000091', NULL, '00000000-0000-0000-0000-000000000006', 'Package boundaries make refactoring much easier.', '2024-03-26 09:20:00+00'),
('20000000-0000-0000-0000-000000000162', '10000000-0000-0000-0000-000000000091', NULL, '00000000-0000-0000-0000-000000000004', 'Domain-driven naming really helps.', '2024-03-26 09:40:00+00'),
('20000000-0000-0000-0000-000000000163', '10000000-0000-0000-0000-000000000092', NULL, '00000000-0000-0000-0000-000000000005', 'Sanitization is essential for UGC.', '2024-03-26 10:20:00+00'),
('20000000-0000-0000-0000-000000000164', '10000000-0000-0000-0000-000000000092', NULL, '00000000-0000-0000-0000-000000000003', 'Moderation workflows matter too.', '2024-03-26 10:40:00+00'),
('20000000-0000-0000-0000-000000000165', '10000000-0000-0000-0000-000000000093', NULL, '00000000-0000-0000-0000-000000000007', 'Preventing errors is better than handling them.', '2024-03-26 11:20:00+00'),
('20000000-0000-0000-0000-000000000166', '10000000-0000-0000-0000-000000000093', NULL, '00000000-0000-0000-0000-000000000004', 'Good validations reduce support load.', '2024-03-26 11:40:00+00'),
('20000000-0000-0000-0000-000000000167', '10000000-0000-0000-0000-000000000094', NULL, '00000000-0000-0000-0000-000000000003', 'Method semantics help caches and clients.', '2024-03-26 12:20:00+00'),
('20000000-0000-0000-0000-000000000168', '10000000-0000-0000-0000-000000000094', NULL, '00000000-0000-0000-0000-000000000006', 'PATCH for partial updates is very convenient.', '2024-03-26 12:40:00+00'),
('20000000-0000-0000-0000-000000000169', '10000000-0000-0000-0000-000000000095', NULL, '00000000-0000-0000-0000-000000000005', 'Deduplication avoids notification spam.', '2024-03-26 13:20:00+00'),
('20000000-0000-0000-0000-000000000170', '10000000-0000-0000-0000-000000000095', NULL, '00000000-0000-0000-0000-000000000007', 'Preferences and quiet hours are important too.', '2024-03-26 13:40:00+00'),
('20000000-0000-0000-0000-000000000171', '10000000-0000-0000-0000-000000000096', NULL, '00000000-0000-0000-0000-000000000004', 'Inventing is the best prediction.', '2024-03-26 14:20:00+00'),
('20000000-0000-0000-0000-000000000172', '10000000-0000-0000-0000-000000000096', NULL, '00000000-0000-0000-0000-000000000006', 'A great quote for builders.', '2024-03-26 14:40:00+00'),
('20000000-0000-0000-0000-000000000173', '10000000-0000-0000-0000-000000000097', NULL, '00000000-0000-0000-0000-000000000007', 'Small improvements keep momentum.', '2024-03-26 15:20:00+00'),
('20000000-0000-0000-0000-000000000174', '10000000-0000-0000-0000-000000000097', NULL, '00000000-0000-0000-0000-000000000005', 'Weekly shipping is a good cadence.', '2024-03-26 15:40:00+00'),
('20000000-0000-0000-0000-000000000175', '10000000-0000-0000-0000-000000000098', NULL, '00000000-0000-0000-0000-000000000003', 'Preventing errors is the real UX.', '2024-03-26 16:20:00+00'),
('20000000-0000-0000-0000-000000000176', '10000000-0000-0000-0000-000000000098', NULL, '00000000-0000-0000-0000-000000000004', 'Agree. Clear defaults help too.', '2024-03-26 16:40:00+00'),
('20000000-0000-0000-0000-000000000177', '10000000-0000-0000-0000-000000000099', NULL, '00000000-0000-0000-0000-000000000006', 'Naming consistency reduces cognitive load.', '2024-03-26 17:20:00+00'),
('20000000-0000-0000-0000-000000000178', '10000000-0000-0000-0000-000000000099', NULL, '00000000-0000-0000-0000-000000000005', 'It also makes search and docs better.', '2024-03-26 17:40:00+00'),
('20000000-0000-0000-0000-000000000179', '10000000-0000-0000-0000-000000000100', NULL, '00000000-0000-0000-0000-000000000004', 'Reliability is built from defaults and limits.', '2024-03-26 18:20:00+00'),
('20000000-0000-0000-0000-000000000180', '10000000-0000-0000-0000-000000000100', '20000000-0000-0000-0000-000000000179', '00000000-0000-0000-0000-000000000002', 'Timeouts and backpressure are the foundation.', '2024-03-26 18:35:00+00'),
('20000000-0000-0000-0000-000000000181', '10000000-0000-0000-0000-000000000100', NULL, '00000000-0000-0000-0000-000000000007', 'Good defaults prevent incidents.', '2024-03-26 18:45:00+00');

-- PUBLICATION LIKES
INSERT INTO publication_likes (id, user_id, publication_id, created_at) VALUES
-- Likes for publication 1
('30000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', '2024-01-20 11:30:00+00'),
('30000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000001', '2024-01-20 12:00:00+00'),
('30000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000001', '2024-01-20 13:00:00+00'),
('30000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000001', '2024-01-21 09:00:00+00'),

-- Likes for publication 6
('30000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000006', '2024-02-05 14:30:00+00'),
('30000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000006', '2024-02-05 15:00:00+00'),
('30000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000006', '2024-02-05 16:00:00+00'),
('30000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000006', '2024-02-06 10:00:00+00'),

-- Likes for publication 13
('30000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000013', '2024-01-30 11:30:00+00'),
('30000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000013', '2024-01-30 12:00:00+00'),
('30000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000013', '2024-01-30 13:00:00+00'),
('30000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000013', '2024-01-31 09:00:00+00'),

-- Likes for publication 7
('30000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000007', '2024-02-10 11:00:00+00'),
('30000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000007', '2024-02-10 12:00:00+00'),

-- Likes for publication 14
('30000000-0000-0000-0000-000000000015', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000014', '2024-02-12 12:30:00+00'),
('30000000-0000-0000-0000-000000000016', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000014', '2024-02-12 13:00:00+00');

-- Additional PUBLICATION LIKES for publications 21-100 (moderate: 3-5 likes per publication)
WITH pubs AS (
  SELECT generate_series(21, 100) AS pub_num
),
like_rows AS (
  SELECT
    p.pub_num,
    u.user_id,
    row_number() OVER (ORDER BY p.pub_num, u.user_id::text) AS seq
  FROM pubs p
  CROSS JOIN LATERAL (
    SELECT unnest(ARRAY[
      '00000000-0000-0000-0000-000000000004'::uuid,
      '00000000-0000-0000-0000-000000000005'::uuid,
      '00000000-0000-0000-0000-000000000006'::uuid
    ]) AS user_id
    UNION ALL SELECT '00000000-0000-0000-0000-000000000007'::uuid WHERE (p.pub_num % 4) = 0
    UNION ALL SELECT '00000000-0000-0000-0000-000000000002'::uuid WHERE (p.pub_num % 10) = 0
  ) u
)
INSERT INTO publication_likes (id, user_id, publication_id, created_at)
SELECT
  ('30000000-0000-0000-0000-' || lpad((16 + lr.seq)::text, 12, '0'))::uuid AS id,
  lr.user_id,
  ('10000000-0000-0000-0000-' || lpad(lr.pub_num::text, 12, '0'))::uuid AS publication_id,
  ('2024-03-19 09:05:00+00'::timestamptz
    + (lr.pub_num - 21) * interval '10 minutes'
    + (lr.seq % 5) * interval '2 minutes') AS created_at
FROM like_rows lr
ORDER BY lr.seq;

-- COMMENT LIKES
INSERT INTO comment_likes (id, user_id, comment_id, created_at) VALUES
('40000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000001', '2024-01-20 11:30:00+00'),
('40000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000006', '20000000-0000-0000-0000-000000000001', '2024-01-20 12:00:00+00'),
('40000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000004', '2024-02-05 14:30:00+00'),
('40000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000005', '20000000-0000-0000-0000-000000000004', '2024-02-05 15:00:00+00'),
('40000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000003', '20000000-0000-0000-0000-000000000008', '2024-01-30 11:30:00+00'),
('40000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000004', '20000000-0000-0000-0000-000000000008', '2024-01-30 12:00:00+00');

-- SAVED ITEMS
INSERT INTO saved_items (id, user_id, publication_id, added_at, note) VALUES
('50000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', '2024-01-21 10:00:00+00', 'Inspirational quote'),
('50000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000013', '2024-01-31 10:00:00+00', 'Reference for microservices'),
('50000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000006', '2024-02-06 11:00:00+00', NULL),
('50000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000014', '2024-02-13 10:00:00+00', 'API design best practices'),
('50000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000013', '2024-02-01 09:00:00+00', 'Great article'),
('50000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000001', '2024-01-22 08:00:00+00', NULL);

-- USER FOLLOWS
INSERT INTO user_follows (id, follower_id, following_id, created_at) VALUES
-- Super admin follows expert and creator
('60000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000002', '2024-01-16 10:00:00+00'),
('60000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000003', '2024-02-02 10:00:00+00'),

-- Expert follows creator
('60000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000003', '2024-02-05 11:00:00+00'),

-- User1 follows expert, creator, and user2
('60000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000002', '2024-02-11 10:00:00+00'),
('60000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000003', '2024-02-11 11:00:00+00'),
('60000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000005', '2024-02-13 10:00:00+00'),

-- User2 follows expert
('60000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000002', '2024-02-13 11:00:00+00'),

-- User3 follows user1 and user2
('60000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000004', '2024-03-02 10:00:00+00'),
('60000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000005', '2024-03-02 11:00:00+00'),

-- Reader follows expert, creator, and user1
('60000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000002', '2024-03-16 10:00:00+00'),
('60000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000003', '2024-03-16 11:00:00+00'),
('60000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000004', '2024-03-16 12:00:00+00');

-- TAGS
INSERT INTO tags (id, name, description, usage_count, created_at) VALUES
('70000000-0000-0000-0000-000000000001', 'programming', 'Programming and software development', 8, '2024-01-15 10:00:00+00'),
('70000000-0000-0000-0000-000000000002', 'microservices', 'Microservices architecture', 5, '2024-01-20 10:00:00+00'),
('70000000-0000-0000-0000-000000000003', 'api-design', 'API design and best practices', 4, '2024-01-25 10:00:00+00'),
('70000000-0000-0000-0000-000000000004', 'inspiration', 'Inspirational quotes and thoughts', 3, '2024-02-01 10:00:00+00'),
('70000000-0000-0000-0000-000000000005', 'web-development', 'Web development topics', 6, '2024-02-05 10:00:00+00'),
('70000000-0000-0000-0000-000000000006', 'database', 'Database optimization and design', 3, '2024-02-10 10:00:00+00'),
('70000000-0000-0000-0000-000000000007', 'security', 'Security best practices', 2, '2024-02-15 10:00:00+00'),
('70000000-0000-0000-0000-000000000008', 'go', 'Go programming language', 2, '2024-03-01 10:00:00+00'),
('70000000-0000-0000-0000-000000000009', 'ai', 'Artificial intelligence', 2, '2024-03-05 10:00:00+00'),
('70000000-0000-0000-0000-000000000010', 'distributed-systems', 'Distributed systems', 1, '2024-03-10 10:00:00+00');

-- PUBLICATION TAGS
INSERT INTO publication_tags (publication_id, tag_id, created_at) VALUES
-- Tags for articles
('10000000-0000-0000-0000-000000000013', '70000000-0000-0000-0000-000000000001', '2024-01-30 10:00:00+00'),
('10000000-0000-0000-0000-000000000013', '70000000-0000-0000-0000-000000000002', '2024-01-30 10:00:00+00'),
('10000000-0000-0000-0000-000000000013', '70000000-0000-0000-0000-000000000010', '2024-01-30 10:00:00+00'),

('10000000-0000-0000-0000-000000000014', '70000000-0000-0000-0000-000000000001', '2024-02-12 11:30:00+00'),
('10000000-0000-0000-0000-000000000014', '70000000-0000-0000-0000-000000000003', '2024-02-12 11:30:00+00'),

('10000000-0000-0000-0000-000000000015', '70000000-0000-0000-0000-000000000001', '2024-02-25 13:45:00+00'),
('10000000-0000-0000-0000-000000000015', '70000000-0000-0000-0000-000000000007', '2024-02-25 13:45:00+00'),
('10000000-0000-0000-0000-000000000015', '70000000-0000-0000-0000-000000000005', '2024-02-25 13:45:00+00'),

('10000000-0000-0000-0000-000000000016', '70000000-0000-0000-0000-000000000001', '2024-03-01 09:00:00+00'),
('10000000-0000-0000-0000-000000000016', '70000000-0000-0000-0000-000000000008', '2024-03-01 09:00:00+00'),

('10000000-0000-0000-0000-000000000017', '70000000-0000-0000-0000-000000000006', '2024-03-08 15:20:00+00'),
('10000000-0000-0000-0000-000000000017', '70000000-0000-0000-0000-000000000001', '2024-03-08 15:20:00+00'),

('10000000-0000-0000-0000-000000000020', '70000000-0000-0000-0000-000000000001', '2024-03-18 14:30:00+00'),
('10000000-0000-0000-0000-000000000020', '70000000-0000-0000-0000-000000000010', '2024-03-18 14:30:00+00'),

-- Tags for posts
('10000000-0000-0000-0000-000000000006', '70000000-0000-0000-0000-000000000001', '2024-02-05 13:00:00+00'),
('10000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000001', '2024-02-10 10:30:00+00'),
('10000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000009', '2024-02-10 10:30:00+00'),
('10000000-0000-0000-0000-000000000007', '70000000-0000-0000-0000-000000000005', '2024-02-10 10:30:00+00'),
('10000000-0000-0000-0000-000000000008', '70000000-0000-0000-0000-000000000001', '2024-02-18 15:00:00+00'),
('10000000-0000-0000-0000-000000000008', '70000000-0000-0000-0000-000000000006', '2024-02-18 15:00:00+00'),

-- Tags for quotes
('10000000-0000-0000-0000-000000000001', '70000000-0000-0000-0000-000000000004', '2024-01-20 10:00:00+00'),
('10000000-0000-0000-0000-000000000002', '70000000-0000-0000-0000-000000000004', '2024-01-25 14:30:00+00'),
('10000000-0000-0000-0000-000000000003', '70000000-0000-0000-0000-000000000004', '2024-02-01 09:15:00+00');

-- Additional PUBLICATION TAGS for publications 21-100 (1-3 tags each)
WITH pubs AS (
  SELECT generate_series(21, 100) AS pub_num
),
pub_tags AS (
  -- Always assign 'programming'
  SELECT pub_num, '70000000-0000-0000-0000-000000000001'::uuid AS tag_id FROM pubs
  UNION ALL
  -- Add a second tag based on modulus for variety
  SELECT pub_num,
         CASE
           WHEN (pub_num % 3) = 0 THEN '70000000-0000-0000-0000-000000000002'::uuid -- microservices
           WHEN (pub_num % 3) = 1 THEN '70000000-0000-0000-0000-000000000003'::uuid -- api-design
           ELSE '70000000-0000-0000-0000-000000000006'::uuid -- database
         END AS tag_id
  FROM pubs
  UNION ALL
  -- Add a third tag only for some publications
  SELECT pub_num,
         CASE
           WHEN (pub_num % 2) = 0 THEN '70000000-0000-0000-0000-000000000008'::uuid -- go
           ELSE '70000000-0000-0000-0000-000000000009'::uuid -- ai
         END AS tag_id
  FROM pubs
  WHERE (pub_num % 4) = 0
)
INSERT INTO publication_tags (publication_id, tag_id, created_at)
SELECT
  ('10000000-0000-0000-0000-' || lpad(pub_num::text, 12, '0'))::uuid AS publication_id,
  tag_id,
  ('2024-03-19 09:00:00+00'::timestamptz + (pub_num - 21) * interval '5 minutes') AS created_at
FROM pub_tags
ORDER BY pub_num, tag_id::text;

-- NOTIFICATIONS
INSERT INTO notifications (id, user_id, type, title, message, data, is_read, created_at) VALUES
-- Notifications for user1
('80000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', 'like', 'New like on your publication', 'testuser2 liked your publication', '{"publication_id": "10000000-0000-0000-0000-000000000008"}', false, '2024-02-18 16:00:00+00'),
('80000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', 'comment', 'New comment on your publication', 'regular_user commented on your publication', '{"publication_id": "10000000-0000-0000-0000-000000000008", "comment_id": "20000000-0000-0000-0000-000000000012"}', true, '2024-02-10 12:30:00+00'),
('80000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000004', 'follow', 'New follower', 'regular_user started following you', '{"follower_id": "00000000-0000-0000-0000-000000000006"}', false, '2024-03-02 10:30:00+00'),

-- Notifications for expert
('80000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000002', 'comment', 'New comment on your publication', 'content_creator commented on your article', '{"publication_id": "10000000-0000-0000-0000-000000000013", "comment_id": "20000000-0000-0000-0000-000000000008"}', true, '2024-01-30 11:30:00+00'),
('80000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000002', 'like', 'New like on your publication', 'testuser1 liked your publication', '{"publication_id": "10000000-0000-0000-0000-000000000013"}', true, '2024-01-30 12:00:00+00'),
('80000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000002', 'follow', 'New follower', 'testuser1 started following you', '{"follower_id": "00000000-0000-0000-0000-000000000004"}', true, '2024-02-11 10:30:00+00'),

-- Notifications for creator
('80000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000003', 'like', 'New like on your publication', 'testuser1 liked your post', '{"publication_id": "10000000-0000-0000-0000-000000000007"}', false, '2024-02-10 11:30:00+00'),
('80000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000003', 'comment', 'New comment on your publication', 'testuser1 commented on your post', '{"publication_id": "10000000-0000-0000-0000-000000000007", "comment_id": "20000000-0000-0000-0000-000000000012"}', true, '2024-02-10 12:00:00+00'),

-- Notifications for user2
('80000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000005', 'follow', 'New follower', 'regular_user started following you', '{"follower_id": "00000000-0000-0000-0000-000000000006"}', false, '2024-03-02 11:30:00+00');

-- PUBLICATION VIEWS
INSERT INTO publications_views (view_uuid, user_id, publication_id, viewed_at) VALUES
-- Views for publication 1
('90000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000001', '2024-01-20 11:00:00+00'),
('90000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000001', '2024-01-20 12:00:00+00'),
('90000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000001', '2024-01-20 13:00:00+00'),
('90000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000001', '2024-01-21 09:00:00+00'),

-- Views for publication 6
('90000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000006', '2024-02-05 14:00:00+00'),
('90000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000006', '2024-02-05 15:00:00+00'),
('90000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000006', '2024-02-05 16:00:00+00'),

-- Views for publication 13
('90000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000013', '2024-01-30 11:00:00+00'),
('90000000-0000-0000-0000-000000000009', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000013', '2024-01-30 12:00:00+00'),
('90000000-0000-0000-0000-000000000010', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000013', '2024-01-30 13:00:00+00'),
('90000000-0000-0000-0000-000000000011', '00000000-0000-0000-0000-000000000006', '10000000-0000-0000-0000-000000000013', '2024-01-31 09:00:00+00'),

-- Views for publication 7
('90000000-0000-0000-0000-000000000012', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000007', '2024-02-10 11:00:00+00'),
('90000000-0000-0000-0000-000000000013', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000007', '2024-02-10 12:00:00+00'),

-- Views for publication 14
('90000000-0000-0000-0000-000000000014', '00000000-0000-0000-0000-000000000003', '10000000-0000-0000-0000-000000000014', '2024-02-12 12:00:00+00'),
('90000000-0000-0000-0000-000000000015', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000014', '2024-02-12 13:00:00+00');

-- Additional PUBLICATION VIEWS for publications 21-100 (1-3 views each; unique (user_id, publication_id, viewed_at))
WITH pubs AS (
  SELECT generate_series(21, 100) AS pub_num
),
view_rows AS (
  SELECT
    p.pub_num,
    u.user_id,
    row_number() OVER (ORDER BY p.pub_num, u.user_id::text) AS seq
  FROM pubs p
  CROSS JOIN LATERAL (
    -- Always: users 4 and 5 view
    SELECT unnest(ARRAY[
      '00000000-0000-0000-0000-000000000004'::uuid,
      '00000000-0000-0000-0000-000000000005'::uuid
    ]) AS user_id
    -- Sometimes: user 6 also views
    UNION ALL SELECT '00000000-0000-0000-0000-000000000006'::uuid WHERE (p.pub_num % 3) = 0
  ) u
)
INSERT INTO publications_views (view_uuid, user_id, publication_id, viewed_at)
SELECT
  ('90000000-0000-0000-0000-' || lpad((15 + vr.seq)::text, 12, '0'))::uuid AS view_uuid,
  vr.user_id,
  ('10000000-0000-0000-0000-' || lpad(vr.pub_num::text, 12, '0'))::uuid AS publication_id,
  ('2024-03-19 09:10:00+00'::timestamptz
    + (vr.pub_num - 21) * interval '7 minutes'
    + (vr.seq % 7) * interval '1 minute') AS viewed_at
FROM view_rows vr
ORDER BY vr.seq;

-- RECOMMENDATIONS
INSERT INTO recommendations (id, user_id, publication_id, algorithm, reason, rank, created_at, hidden) VALUES
-- Recommendations for user1
('a0000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000013', 'collaborative', 'Based on your interests in programming', 0, '2024-02-01 10:00:00+00', false),
('a0000000-0000-0000-0000-000000000002', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000014', 'collaborative', 'Similar users liked this', 1, '2024-02-01 10:00:00+00', false),
('a0000000-0000-0000-0000-000000000003', '00000000-0000-0000-0000-000000000004', '10000000-0000-0000-0000-000000000015', 'content-based', 'Matches your reading history', 2, '2024-02-01 10:00:00+00', false),

-- Recommendations for user2
('a0000000-0000-0000-0000-000000000004', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000013', 'collaborative', 'Based on your interests', 0, '2024-02-13 10:00:00+00', false),
('a0000000-0000-0000-0000-000000000005', '00000000-0000-0000-0000-000000000005', '10000000-0000-0000-0000-000000000017', 'content-based', 'Popular in your network', 1, '2024-02-13 10:00:00+00', false),

-- Recommendations for reader
('a0000000-0000-0000-0000-000000000006', '00000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000001', 'popular', 'Popular quote', 0, '2024-03-16 10:00:00+00', false),
('a0000000-0000-0000-0000-000000000007', '00000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000006', 'popular', 'Trending post', 1, '2024-03-16 10:00:00+00', false),
('a0000000-0000-0000-0000-000000000008', '00000000-0000-0000-0000-000000000007', '10000000-0000-0000-0000-000000000013', 'collaborative', 'Based on followed users', 2, '2024-03-16 10:00:00+00', false);

-- Additional RECOMMENDATIONS for new publications 21-100
WITH pub_set AS (
  SELECT generate_series(21, 100) AS pub_num
),
user_set AS (
  -- Keep it focused on the same test users used above
  SELECT unnest(ARRAY[
    '00000000-0000-0000-0000-000000000004'::uuid, -- testuser1
    '00000000-0000-0000-0000-000000000005'::uuid, -- testuser2
    '00000000-0000-0000-0000-000000000007'::uuid  -- reader_only
  ]) AS user_id
),
ranked AS (
  SELECT
    u.user_id,
    p.pub_num,
    row_number() OVER (PARTITION BY u.user_id ORDER BY p.pub_num) - 1 AS rank,
    row_number() OVER (ORDER BY u.user_id::text, p.pub_num) AS seq
  FROM user_set u
  JOIN pub_set p ON ((p.pub_num + (CASE WHEN u.user_id = '00000000-0000-0000-0000-000000000004'::uuid THEN 0
                                      WHEN u.user_id = '00000000-0000-0000-0000-000000000005'::uuid THEN 1
                                      ELSE 2 END)) % 3) = 0
)
INSERT INTO recommendations (id, user_id, publication_id, algorithm, reason, rank, created_at, hidden)
SELECT
  ('a0000000-0000-0000-0000-' || lpad((8 + r.seq)::text, 12, '0'))::uuid AS id,
  r.user_id,
  ('10000000-0000-0000-0000-' || lpad(r.pub_num::text, 12, '0'))::uuid AS publication_id,
  CASE WHEN (r.pub_num % 2) = 0 THEN 'collaborative' ELSE 'content-based' END AS algorithm,
  CASE WHEN (r.pub_num % 2) = 0 THEN 'Similar users liked this' ELSE 'Matches your reading history' END AS reason,
  r.rank,
  ('2024-03-19 10:00:00+00'::timestamptz + r.rank * interval '3 minutes') AS created_at,
  false AS hidden
FROM ranked r
ORDER BY r.seq;

COMMIT;

