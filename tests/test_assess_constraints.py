"""Tests for assess_constraints module."""

import logging
from mock import Mock, patch
import StringIO
import os
from contextlib import contextmanager

from assess_constraints import (
    append_constraint,
    make_constraints,
    assess_virt_type_constraints,
    assess_instance_type_constraints,
    deploy_charm_constraint,
    parse_args,
    main,
    INSTANCE_TYPES,
    )
from tests import (
    parse_error,
    TestCase,
    )
from tests.test_jujupy import fake_juju_client
from utility import temp_dir


class TestParseArgs(TestCase):

    def test_common_args(self):
        args = parse_args(["an-env", "/bin/juju", "/tmp/logs", "an-env-mod"])
        self.assertEqual("an-env", args.env)
        self.assertEqual("/bin/juju", args.juju_bin)
        self.assertEqual("/tmp/logs", args.logs)
        self.assertEqual("an-env-mod", args.temp_env_name)
        self.assertEqual(False, args.debug)

    def test_help(self):
        fake_stdout = StringIO.StringIO()
        with parse_error(self) as fake_stderr:
            with patch("sys.stdout", fake_stdout):
                parse_args(["--help"])
        self.assertEqual("", fake_stderr.getvalue())


class TestMakeConstraints(TestCase):

    def test_append_constraint_none(self):
        args = []
        append_constraint(args, 'name', None)
        self.assertEqual([], args)

    def test_append_constraint_string(self):
        args = ['inital=True']
        append_constraint(args, 'name', 'value')
        self.assertEqual(['inital=True', 'name=value'], args)

    def test_make_constraints_empty(self):
        constraints = make_constraints()
        self.assertEqual('', constraints)

    def test_make_constraints(self):
        constraints = make_constraints(memory='20GB', virt_type='test')
        if 'm' == constraints[0]:
            self.assertEqual('mem=20GB virt-type=test', constraints)
        else:
            self.assertEqual('virt-type=test mem=20GB', constraints)


class TestMain(TestCase):

    def test_main(self):
        argv = ["an-env", "/bin/juju", "/tmp/logs", "an-env-mod", "--verbose"]
        client = Mock(spec=["is_jes_enabled"])
        with patch("assess_constraints.configure_logging",
                   autospec=True) as mock_cl:
            with patch("assess_constraints.BootstrapManager.booted_context",
                       autospec=True) as mock_bc:
                with patch('deploy_stack.client_from_config',
                           return_value=client) as mock_cfc:
                    with patch("assess_constraints.assess_constraints",
                               autospec=True) as mock_assess:
                        main(argv)
        mock_cl.assert_called_once_with(logging.DEBUG)
        mock_cfc.assert_called_once_with('an-env', "/bin/juju", debug=False)
        self.assertEqual(mock_bc.call_count, 1)
        mock_assess.assert_called_once_with(client, False)


class TestAssess(TestCase):

    @contextmanager
    def prepare_deploy_mock(self):
        with patch('jujupy.EnvJujuClient.deploy',
                   autospec=True) as deploy_mock:
            yield deploy_mock

    def test_virt_type_constraints_with_kvm(self):
        # Using fake_client means that deploy and get_status have plausible
        # results.  Wrapping it in a Mock causes every call to be recorded, and
        # allows assertions to be made about calls.  Mocks and the fake client
        # can also be used separately.
        fake_client = Mock(wraps=fake_juju_client())
        assert_constraints_calls = ["virt-type=lxd", "virt-type=kvm"]
        fake_client.bootstrap()
        deploy = patch('jujupy.EnvJujuClient.deploy',
                       autospec=True)
        with deploy as deploy_mock:
            assess_virt_type_constraints(fake_client, True)
        constraints_calls = [
            call[1]["constraints"] for call in
            deploy_mock.call_args_list]
        self.assertEqual(constraints_calls, assert_constraints_calls)

    def test_virt_type_constraints_without_kvm(self):
        # Using fake_client means that deploy and get_status have plausible
        # results.  Wrapping it in a Mock causes every call to be recorded, and
        # allows assertions to be made about calls.  Mocks and the fake client
        # can also be used separately.
        fake_client = Mock(wraps=fake_juju_client())
        assert_constraints_calls = ["virt-type=lxd"]
        fake_client.bootstrap()
        deploy = patch('jujupy.EnvJujuClient.deploy',
                       autospec=True)
        with deploy as deploy_mock:
            assess_virt_type_constraints(fake_client, False)
        constraints_calls = [
            call[1]["constraints"] for call in
            deploy_mock.call_args_list]
        self.assertEqual(constraints_calls, assert_constraints_calls)

    # TODO Get these two working:
    def do_not_test_instance_type_constraints(self):
        fake_client = Mock(wraps=fake_juju_client())
        #fake_client.env.config['type'] = 'foo'
        INSTANCE_TYPES['foo'] = ['bar' , 'baz']
        with patch('assess_constraints.assess_instance_type',
                   autospec=True) as assess_mock:
            assess_instance_type_constraints(fake_client)
        calls = assess_mock.call_args_list
        self.assertEqual(len(INSTANCE_TYPES[provider]), len(calls))
        for instance_type in INSTANCE_TYPES[provider]:
            pass
        del INSTANCE_TYPES['foo']

    def test_instance_type_constraints_missing(self):
        fake_client = Mock(wraps=fake_juju_client())


class TestDeploy(TestCase):

    def test_deploy_charm_constraint(self):
        fake_client = Mock(wraps=fake_juju_client())
        charm_name = 'test-constraint'
        charm_series = 'xenial'
        constraints = 'mem=10GB'
        with temp_dir() as charm_dir:
            with patch('assess_constraints.deploy_constraint',
                       autospec=True) as deploy_mock:
                deploy_charm_constraint(fake_client, charm_name, charm_series,
                                        charm_dir, constraints)
        charm = os.path.join(charm_dir, charm_series, charm_name)
        deploy_mock.assert_called_once_with(fake_client, charm, charm_series,
                                            charm_dir, constraints)
