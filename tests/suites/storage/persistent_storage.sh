# This function tests charm deployment with a single filesystem storage unit and a single persistent block device
# storage unit.
#  Steps taken to the test:
#       - Deploy dummy-storage charm with a single block storage unit (ebs)
#         and a single filesystem storage unit (rootfs).
#       - Check charm status once the deployment is done.
#           > Application status should be active.
#       - Check charm storage units once the deployment is done.
#           > Total number of storage units should be 2.
#           > Name of storage units should be in align with charm config.
#           > Properties of storage units should be as defined.
#               - Storage Type, Persistent Setting and Pool.
run_persistent_storage() {
	echo

	model_name="persistent-storage"
	file="${TEST_DIR}/test-${model_name}.log"
	ensure "${model_name}" "${file}"
	# dummy-storage is going to be deployed with 1 ebs block storage unit
	# and 1 rootfs filesystem storage unit.
	echo "dummy-storage is going to be deployed with 1 ebs block storage unit and 1 rootfs filesystem storage unit"
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage --storage single-blk=ebs \
		--storage single-fs=rootfs
	echo "Checking current status of app dummy-storage."
	# wait for current application-status to be active
	wait_for "dummy-storage" "$(active_condition "dummy-storage" 0)"
	# wait for current workload-status to be active
	wait_for "active" "$(workload_status "dummy-storage" 0).current"

	echo "Checking total number of storage unit(s)."
	juju storage --format json | jq '.storage | keys | length' | check 2
	echo "Checking names of storage unit(s)."
	juju storage --format json | jq '.storage | keys | .[0]' | check "single-blk/0"
	juju storage --format json | jq '.storage | keys | .[1]' | check "single-fs/1"
	echo "Check name and total number of storage unit: PASSED."

	#
	# check type, persistent setting and pool of single block storage unit
	#
	# storage type
	juju storage --format json | jq '.storage | ."single-blk/0" | .kind' | check "block"
	# persistent setting
	echo "Checking persistent setting of single block storage unit"
	juju storage --format json | jq '.storage | ."single-blk/0" | .persistent' | check true
	# pool storage
	juju storage --format json | jq '.volumes | ."0" | .storage' | check "single-blk/0"
	# pool setting
	juju storage --format json | jq '.volumes | ."0" | .pool' | check "ebs"
	# storage status
	juju storage --format json | jq '.volumes | ."0" | .status | .current' | check "attached"
	echo "Check properties of single block device storage unit: PASSED."

	#
	# check type, persistent setting and pool of single filesystem storage unit
	#
	juju storage --format json | jq '.storage | ."single-fs/1" | .kind' | check "filesystem"
	# persistent setting
	echo "Checking persistent setting of single filesystem storage unit"
	juju storage --format json | jq '.storage | ."single-fs/1" | .persistent' | check false
	# pool storage
	juju storage --format json | jq '.volumes | ."0" | .storage' | check "single-blk/0"
	# pool setting
	juju storage --format json | jq '.volumes | ."0" | .pool' | check "ebs"
	# storage status
	juju storage --format json | jq '.volumes | ."0" | .status | .current' | check "attached"
	echo "Check properties of single filesystem storage unit: PASSED."

	# assert charm removal message for single block and filesystem storage
	removal_msg=$(juju remove-application dummy-storage 2>&1)
	echo "${removal_msg}" | sed -sn 2p | sed 's/^-//' | check "will remove storage single-fs/1"
	echo "${removal_msg}" | sed -sn 3p | sed 's/^-//' | check "will detach storage single-blk/0"
	#
	# wait until an update of storage status occurred. Due to the asynchronous nature of Juju,
	# the status of storage may change time to time after a Juju CLI issued, in this test only
	# the existence of storage id is the point of interest.
	#
	wait_for "{}" ".applications" # we use this wait_for command as an indicator that the storage
	# status has changed and now we can check for the storage status and assert that indeed the
	# single filesystem storage unit has been removed successfully.
	juju storage --format json | jq '.storage | has("single-fs/1")' | check false

	echo "Checking total number of storage unit(s)."
	juju storage --format json | jq '.storage | keys | length' | check 1
	echo "Check for existence of single block storage"
	juju storage --format json | jq '.storage | has("single-blk/0")' | check true
	echo "single-blk/0 found in storage list."

	echo "Check for existence of single-blk/0 persistent storage after remove-application"
	juju storage --format json | jq '.volumes | ."0" | .storage' | check "single-blk/0"
	# storage status
	juju storage --format json | jq '.volumes | ."0" | .status | .current' | check "detached"
	echo "Check status of persistent storage single-blk/0 after remove-application: PASSED"

	# Deploy charm with an existing detached storage
	juju deploy -m "${model_name}" ./testcharms/charms/dummy-storage --attach-storage single-blk/0
	echo "Checking current status of app dummy-storage."
	# wait for current application-status to be active
	wait_for "dummy-storage" "$(active_condition "dummy-storage" 0)"
	# wait for current workload-status to be active
	wait_for "active" "$(workload_status "dummy-storage" 1).current"
	# assert storage unit count
	juju storage --format json | jq '.storage | keys | length' | check 1
	echo "Checking existence of single block device storage single-blk/0."
	juju storage --format json | jq '.storage | has("single-blk/0")' | check true
	# assert persistent setting
	echo "Checking persistent setting of storage unit single-blk/0."
	juju storage --format json | jq '.storage | ."single-blk/0" | .persistent' | check true
	# assert storage status
	echo "Checking the status of storage single-blk/0 in volumes."
	juju storage --format json | jq '.volumes | ."0" | .status | .current' | check "attached"

	# remove charm
	juju remove-application dummy-storage
	# wait for application to be removed
	wait_for "{}" ".applications"
	# persistent storage should remain after remove-application
	juju storage --format json | jq '.storage | has("single-blk/0")' | check true
	# remove storage
	juju remove-storage single-blk/0
	# wait until an update of storage status occurred. Due to the asynchronous nature of Juju,
	# the status of storage may change time to time after a Juju CLI issued.
	# Checking removal of single block device storage {}.
	# Bug 1708340 https://bugs.launchpad.net/juju/+bug/1708340
	# juju storage --format json output will be empty if no storage exists. So we use juju status --format json
	# and check for the storage section.
	wait_for "{}" ".storage"
	destroy_model "persistent_storage"
}

test_persistent_storage() {
	if [ "$(skip 'test-persistent-storage')" ]; then
		echo "==> TEST SKIPPED: persistent storage tests"
		return
	fi

	(
		set_verbosity

		cd .. || exit

		run "run_persistent_storage"
	)
}
