# GitHub Action Setup for Datasette Deployment

## Changes Made

✅ Updated `.github/workflows/datasette-deploy.yml` to:
- Deploy using Vercel CLI instead of `datasette publish vercel`
- Deploy the `datasette-deploy` directory with custom `index.py`
- Connect to Postgres via `POSTGRES_URL` environment variable
- Trigger on push to main/master when datasette-deploy files change

## Required GitHub Secrets

For the GitHub Action to work, you need to add these secrets to your GitHub repository:

### 1. VERCEL_TOKEN
Get from: https://vercel.com/account/tokens

### 2. VERCEL_ORG_ID
Get from Vercel project settings:
- Go to your project in Vercel
- Settings → General
- Copy "Team ID" or "Organization ID"

### 3. VERCEL_PROJECT_ID  
Get from Vercel project settings:
- Go to your datasette project in Vercel
- Settings → General
- Copy "Project ID"

## How to Add Secrets to GitHub

1. Go to your GitHub repository
2. Settings → Secrets and variables → Actions
3. Click "New repository secret"
4. Add each of the three secrets above

## Getting Vercel IDs via CLI

```bash
cd /workspace/datasette-deploy
vercel link
# This creates .vercel/project.json with the IDs

cat .vercel/project.json
```

Copy the `orgId` and `projectId` values.

## Important: Environment Variables in Vercel

The GitHub Action will deploy your code, but **Vercel still needs the `POSTGRES_URL` environment variable set in the dashboard**:

1. Go to Vercel Dashboard
2. Select your datasette project
3. Settings → Environment Variables
4. Ensure `POSTGRES_URL` is set (you confirmed this is already done ✅)

## How the Workflow Works

### Triggers:
- **Manual:** Via workflow_dispatch (Actions tab → Run workflow)
- **Auto:** Push to main/master that changes `datasette-deploy/**` files
- **Scheduled:** Every 6 hours (for data updates)

### Steps:
1. Checkout code
2. Install Vercel CLI
3. Deploy to Vercel using `vercel deploy --prod`
4. Show deployment summary

## Testing the Workflow

### Option 1: Manual Trigger (Recommended for testing)
1. Go to GitHub → Actions tab
2. Select "Deploy Datasette to Vercel"
3. Click "Run workflow"
4. Watch the logs

### Option 2: Push Changes
The workflow will auto-trigger when you push changes to `datasette-deploy/` files.

## What Happens if Secrets are Missing

If secrets are not set, the workflow will:
- ✅ Not fail
- ℹ️ Print helpful message about how to add secrets
- ⏭️ Skip deployment gracefully

## Next Steps

1. ✅ Code is committed (0e87ab2)
2. ⬜ Push to GitHub (system will handle automatically)
3. ⬜ Add GitHub secrets (VERCEL_TOKEN, VERCEL_ORG_ID, VERCEL_PROJECT_ID)
4. ⬜ Test workflow manually from GitHub Actions tab
5. ⬜ Verify deployment works with Postgres

## Summary

- **Code:** ✅ Ready and committed
- **Workflow:** ✅ Updated to use Postgres
- **Postgres URL:** ✅ Confirmed set in Vercel
- **GitHub Secrets:** ⬜ Need to be added by you
- **Ready to push:** ✅ Yes!
