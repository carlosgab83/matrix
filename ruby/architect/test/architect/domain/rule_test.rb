# typed: false
# frozen_string_literal: true

require 'test_helper'

module Architect
  module Domain
    class RuleTest < Minitest::Test
      def test_creates_rule_with_name
        rule = Rule.new(
          name: 'test_rule',
          conditions: [],
          actions: [],
          inclusion: [],
          exclusion: []
        )

        assert_equal rule.name, 'test_rule'
        assert_equal rule.conditions, []
        assert_equal rule.actions, []
      end

      def test_creates_rule_must_fail_if_inclusion_and_exclusion_are_present
        assert_raises(Architect::Domain::Errors::InvalidRuleDefinition) do
          Rule.new(
            name: 'test_rule',
            conditions: [],
            actions: [],
            inclusion: [:any],
            exclusion: [:any]
          )
        end
      end
    end
  end
end
