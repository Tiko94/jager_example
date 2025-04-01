SERVICES=dataservice loadbalancer logservice web
VERSION=$(shell cat VERSION)
ECR_BASE=589295909756.dkr.ecr.us-east-2.amazonaws.com
CHARTS_REPO_NAME=helm-rnd-charts

CHART=jaeger-rd

.PHONY: all clean
all: build-all-images publish-all-images build-chart publish-chart

build-image-%:
	@echo "+ $@"
	docker build -t ${ECR_BASE}/jaeger-rd-$*:${VERSION} -f Dockerfile-$* .

publish-image-%: build-image-%
	@echo "+ $@"
	docker push ${ECR_BASE}/jaeger-rd-$*:${VERSION}

build-all-images: $(addprefix build-image-,$(SERVICES))
publish-all-images: $(addprefix publish-image-,$(SERVICES))

build-chart:
	@echo "+ $@"
	cd helm-charts && \
	helm lint ${CHART} && \
	helm dependency build ${CHART} && \
	helm package --app-version "${VERSION}" --version "${VERSION}" "${CHART}"

publish-chart: build-chart
	@echo "+ $@"
	helm s3 push --force helm-charts/${CHART}-${VERSION}.tgz ${CHARTS_REPO_NAME}

install-microk8s:
	@echo "+ $@"
	helm repo update && \
	helm upgrade \
		--install \
		--values helm-charts/${CHART}-microk8s-values.yaml \
		--version ${VERSION} \
		blog \
		${CHARTS_REPO_NAME}/jaeger-rd

install-eks:
	@echo "+ $@"
	helm repo update && \
	helm upgrade \
		--install \
		--values helm-charts/${CHART}-eks-values.yaml \
		--version ${VERSION} \
		blog \
		${CHARTS_REPO_NAME}/jaeger-rd

uninstall:
	@echo "+ $@"
	helm uninstall blog

clean:
	@echo "+ $@"
	cd helm-charts && \
	rm -fv *.tgz
