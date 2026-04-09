# Hosting the Wivvus API on DigitalOcean App Platform

This guide covers deploying the Wivvus Go API using DigitalOcean App Platform, which deploys directly from your GitHub repository and automatically redeploys on every push to `main`.

---

## Overview

DigitalOcean App Platform:

- Pulls code directly from GitHub
- Builds and runs the Go binary automatically
- Manages TLS certificates and a public HTTPS URL
- Provides a managed PostgreSQL database add-on
- Injects environment variables securely at runtime
- Redeploys automatically on every push to `main`

---

## Prerequisites

- A DigitalOcean account
- The API repository pushed to GitHub (`github.com/Wivvus/api`)
- A Google OAuth client ID and secret
- A Gmail App Password (or other SMTP credentials)

---

## 1. Create the App

1. In the DigitalOcean control panel, go to **App Platform → Create App**.
2. Choose **GitHub** as the source.
3. Authorise DigitalOcean to access your GitHub account if prompted.
4. Select the `Wivvus/api` repository and the `main` branch.
5. Enable **Autodeploy** — this redeploys the app automatically on every push.
6. Click **Next**.

---

## 2. Configure the Service

App Platform will detect the Go module and suggest a Web Service component.

On the component configuration screen:

- **Type:** Web Service
- **Build Command:** leave blank (App Platform detects `go.mod` and builds automatically)
- **Run Command:** `go run cmd/api/main.go` — or, for a faster start, set a build command of `go build -o api-server ./cmd/api` and a run command of `./api-server`
- **HTTP Port:** `8080`

To use the compiled binary (recommended for production):

| Field | Value |
|---|---|
| Build Command | `go build -o api-server ./cmd/api` |
| Run Command | `./api-server` |
| HTTP Port | `8080` |

Click **Next**.

---

## 3. Add a Managed PostgreSQL Database

1. On the **Add Resources** screen, click **Add a Database**.
2. Select **Dev Database** (free, suitable for development) or a paid cluster for production.
3. Name it `wivvus-db`.
4. App Platform will automatically inject the following environment variables into your app:

| Variable | Provided by |
|---|---|
| `DATABASE_URL` | Auto-injected connection string |

> **Note:** The Wivvus API reads individual `PG_HOST`, `PG_PORT`, `PG_USER`, `PG_PASSWORD`, `PG_DB`, and `PG_SSLMODE` variables rather than a single `DATABASE_URL`. After the database is created, find its connection details under the database component and set them individually in the next step.

---

## 4. Set Environment Variables

On the **Environment Variables** screen, add the following. Mark sensitive values as **Encrypted** (the lock icon) so they are stored securely and never shown in logs.

| Variable | Value | Encrypted |
|---|---|---|
| `PG_HOST` | From the database component's connection details | Yes |
| `PG_PORT` | From the database component (usually `25060`) | No |
| `PG_USER` | From the database component | Yes |
| `PG_PASSWORD` | From the database component | Yes |
| `PG_DB` | From the database component | No |
| `PG_SSLMODE` | `require` | No |
| `JWT_SECRET` | A long random string — generate with `openssl rand -hex 32` | Yes |
| `GOOGLE_OAUTH_CLIENT_ID` | Your Google OAuth client ID | No |
| `SMTP_HOST` | `smtp.gmail.com` | No |
| `SMTP_PORT` | `587` | No |
| `SMTP_USER` | Your Gmail address | No |
| `SMTP_PASSWORD` | Your Gmail App Password | Yes |
| `SMTP_FROM` | Your Gmail address | No |
| `APP_URL` | Your frontend URL e.g. `https://wivvus.com` | No |
| `ALLOWED_ORIGINS` | Your frontend URL e.g. `https://wivvus.com` | No |

### Finding database connection details

After adding the database component, click on it in the App Platform dashboard. Under **Connection Details**, choose **App-layer connection string components** to see the individual host, port, user, password, and database name values.

### Generating a JWT secret

Run this locally and paste the output:

```bash
openssl rand -hex 32
```

---

## 5. Choose a Plan and Deploy

1. Select a plan — the **Basic** plan ($5/month) is sufficient to start.
2. Review the summary and click **Create Resources**.

App Platform will clone the repository, build the Go binary, and start the service. The first deploy takes a few minutes.

---

## 6. Get Your API URL

Once deployed, App Platform assigns a public URL like:

```
https://wivvus-api-xxxx.ondigitalocean.app
```

Find it under **App Overview → Live URL**. Copy it — you will need it for the frontend environment config and the `ALLOWED_ORIGINS` variable.

---

## 7. Add a Custom Domain (Optional)

1. In the App settings, go to **Domains → Add Domain**.
2. Enter your domain, e.g. `api.wivvus.com`.
3. Follow the instructions to add a CNAME record in your DNS provider pointing to the App Platform URL.
4. App Platform provisions a Let's Encrypt certificate automatically.

Once the custom domain is active, update `ALLOWED_ORIGINS` to use it.

---

## 8. Create a Google OAuth Client

1. Go to the [Google Cloud Console](https://console.cloud.google.com) and create a project (or select an existing one).
2. Go to **APIs & Services → OAuth consent screen**.
   - Set **User type** to **External**.
   - Fill in the app name, support email, and developer contact email. The other fields can be left blank.
   - Click through to **Save and Continue** on the remaining steps.
3. Go to **APIs & Services → Credentials → Create Credentials → OAuth 2.0 Client ID**.
   - Set **Application type** to **Web application**.
   - Under **Authorised JavaScript origins**, add your frontend URL:
     ```
     https://wivvus.com
     ```
   - Under **Authorised redirect URIs**, add:
     ```
     https://<your-api-url>/auth/google/callback
     ```
4. Click **Create**. Copy the **Client ID** — this goes into `GOOGLE_OAUTH_CLIENT_ID` in step 4. The client secret is not required by the API.

> If you don't yet know your API URL (assigned after deployment in step 6), you can come back and add the redirect URI afterwards. The app won't crash on startup — the redirect URI is only checked by Google at login time.

---

## Ongoing Operations

### Deploying an update

Push to `main` — App Platform detects the push and redeploys automatically.

```bash
git push origin main
```

### Viewing logs

In the App Platform dashboard, go to **Runtime Logs** to see live output from the running service.

### Updating environment variables

Go to **App Settings → App-Level Environment Variables**. Changes trigger a redeploy.

### Scaling

To handle more traffic, go to **App Settings → Components** and increase the instance size or count without any config file changes.
