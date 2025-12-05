# typed: false
# frozen_string_literal: true

module Architect
  module Domain
    class Action
      attr_reader :action_type, :action_value

      def initialize(action_type:, action_value:)
        @action_type = action_type
        @action_value = action_value
        validate_action!
      end

      private

      def validate_action!
        return if action_type && action_value

        raise Architect::Domain::Errors::InvalidRuleDefinition,
              'Invalid action. Actions must have action_type and action_value'
      end
    end
  end
end
