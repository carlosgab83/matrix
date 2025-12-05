# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/errors/invalid_rule_definition'

module Architect
  module Domain
    class Condition
      extend T::Sig

      sig { returns(String) }
      attr_reader :key

      sig { returns(Symbol) }
      attr_reader :operator

      sig { returns(T.untyped) }
      attr_reader :value

      AVAILABLE_OPERATORS = %i[
        equals_to
        not_equals_to
        greater_than
        greater_or_equal
        less_than
        less_or_equal
      ].freeze

      sig { params(key: String, operator: Symbol, value: T.untyped).void }
      def initialize(key:, operator:, value:)
        @key = T.let(key, String)
        @operator = T.let(operator, Symbol)
        @value = T.let(value, T.untyped)
        validate_condition!
      end

      private

      sig { void }
      def validate_condition!
        return if AVAILABLE_OPERATORS.include?(operator)

        raise Architect::Domain::Errors::InvalidRuleDefinition,
              "Invalid condition operator '#{operator}'. Availables are #{AVAILABLE_OPERATORS}"
      end
    end
  end
end
