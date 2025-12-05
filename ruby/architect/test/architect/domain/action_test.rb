# typed: false
# frozen_string_literal: true

require 'test_helper'

module Architect
  module Domain
    class ActionTest < Minitest::Test
      def test_creates_action_with_params
        action = Action.new(
          action_type: 'any_type',
          action_value: 'any_value'
        )

        assert_equal action.action_type, 'any_type'
        assert_equal action.action_value, 'any_value'
      end

      def test_creates_action_must_fail_if_empty
        assert_raises(Architect::Domain::Errors::InvalidRuleDefinition) do
          Action.new(
            action_type: nil,
            action_value: nil
          )
        end
      end
    end
  end
end
