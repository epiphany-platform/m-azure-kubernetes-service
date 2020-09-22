test-default-config:
	#will run default config tests
	@bash tests/tests.sh cleanup
	@bash tests/tests.sh setup
	@bash tests/tests.sh test-default-config-suite $(IMAGE_NAME)
	@bash tests/tests.sh cleanup
	#finished default config tests

test-config-with-variables:
	# tests of config with variables will go here

test-plan:
	# tests of plan operation will go here

test-apply:
	# tests of apply operation will go here

generate-report:
	@bash tests/tests.sh generate_junit_report
