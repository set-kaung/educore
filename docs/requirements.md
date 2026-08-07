Core Requirements to hit Course Objectives

    Infrastructure: Deployed on a hardened Linux VPS (Oracle Cloud or Azure).

    Networking & Deployment: Exposed via Nginx Reverse Proxy with Let's Encrypt SSL. The project must be deployed under a new, distinct URL path (e.g., https://your-domain.com/project). This requires writing custom Nginx location blocks or Docker port mappings that do not break your existing /content (WordPress) or /api (Lab) routes.

    Backend: Node.js (Express) or Go REST API.

    Database: Relational Database (MySQL/PostgreSQL) managed via Prisma ORM with migrations.

    Security & Identity (User Auth): JWT Authentication with Role-Based Access Control (RBAC),

    integrating the University's Microsoft Active Directory (AD) for user authentication (using MSAL, OAuth2, or OIDC).

    Secrets Management (Azure Key Vault): You must NOT use local .env files for production secrets. Instead, your application must authenticate with the centralized Class Azure Key Vault server (credentials will be provided) to securely fetch your database connection strings, JWT secrets, and API keys at runtime.

    External Integration (3rd Party): The backend must connect to and utilize at least one external public API or AI service (e.g., OpenAI, Gemini, Google Maps, SendGrid, Weather API) to augment its business logic.

    Service-to-Service Integration (Peer API): You must partner with another student/team in the class.

        Expose: You must build an endpoint specifically for their backend to consume. This endpoint must be protected by a static API Key (e.g., expecting an x-api-key header) that you generate and issue exclusively to them.

        Consume: Your backend must fetch and utilize data from their API, authenticating your requests using the API Key they issued to you.

    Source Code Management: All code must be hosted in a GitHub repository.

    Automation: Automated deployment script (or Docker Compose).
