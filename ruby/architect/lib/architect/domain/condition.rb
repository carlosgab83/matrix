# typed: false
# frozen_string_literal: true

require 'architect/domain/errors/invalid_rule_definition'

module Architect
  module Domain
    class Condition
      attr_reader :key, :operator, :value

      AVAILABLE_OPERATORS = %i[
        equals_to
        not_equals_to
        greater_than
        greater_or_equal
        less_than
        less_or_equal
      ].freeze

      def initialize(key:, operator:, value:)
        @key = key
        @operator = operator
        @value = value
        validate_condition!
      end

      private

      def validate_condition!
        return if AVAILABLE_OPERATORS.include?(operator)

        raise Architect::Domain::Errors::InvalidRuleDefinition,
              "Invalid condition operator '#{operator}'. Availables are #{AVAILABLE_OPERATORS}"
      end
    end
  end
end
