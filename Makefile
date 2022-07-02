build-frontend:
	./frontend/bin/build.sh

deploy-backend:
	./backend/bin/deploy.sh

deploy-backend-2:
	./bin/deploy-lambda/run.sh backend/lambdas/db/migrations

deploy-frontend:
	./frontend/bin/deploy.sh
