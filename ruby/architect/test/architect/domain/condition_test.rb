# typed: false
# frozen_string_literal: true

require 'test_helper'

module Architect
  module Domain
    class ConditionTest < Minitest::Test
      def test_creates_condition_with_params
        condition = Condition.new(
          key: 'test_key',
          operator: :equals_to,
          value: 'any'
        )

        assert_equal condition.key, 'test_key'
        assert_equal condition.operator, :equals_to
        assert_equal condition.value, 'any'
      end

      def test_creates_condition_must_fail_if_unavailable_operator
        assert_raises(Architect::Domain::Errors::InvalidRuleDefinition) do
          Condition.new(
            key: 'test_key',
            operator: :_any,
            value: 'any'
          )
        end
      end
    end
  end
end
