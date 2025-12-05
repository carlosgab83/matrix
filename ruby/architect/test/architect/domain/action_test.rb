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
    end
  end
end
