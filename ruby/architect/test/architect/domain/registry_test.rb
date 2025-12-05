# typed: false
# frozen_string_literal: true

require 'test_helper'
require 'architect/shared/app_logger'

module Architect
  module Domain
    class RegistryTest < Minitest::Test
      def setup
        Registry.clear!
        @fixtures_path = File.expand_path('../../fixtures/rules', __dir__)
      end

      def teardown
        Registry.clear!
      end

      def test_creates_rules_from_registry
        Registry.reload_from_dir!(@fixtures_path)

        assert_equal 2, Registry.rules.length

        rule_names = Registry.rules.map(&:name)
        assert_includes rule_names, 'fixture_rule_1'
        assert_includes rule_names, 'fixture_rule_2'
      end
    end
  end
end
