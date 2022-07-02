build-frontend:
	./frontend/bin/build.sh

deploy-backend:
	./bin/deploy-lambda/run.sh backend/lambdas/jobs/pollschedules
	./bin/deploy-lambda/run.sh backend/lambdas/apis/devices
	./bin/deploy-lambda/run.sh backend/lambdas/apis/units
	./bin/deploy-lambda/run.sh backend/lambdas/apis/users
	./bin/deploy-lambda/run.sh backend/lambdas/apis/signin
	./bin/deploy-lambda/run.sh backend/lambdas/apis/properties
	./bin/deploy-lambda/run.sh backend/lambdas/db/migrations

deploy-frontend:
	./frontend/bin/deploy.sh
