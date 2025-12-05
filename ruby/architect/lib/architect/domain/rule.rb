# typed: false
# frozen_string_literal: true

require 'architect/domain/errors/invalid_rule_definition'

module Architect
  module Domain
    class Rule
      attr_reader :name, :conditions, :actions, :inclusion, :exclusion

      def initialize(name:, conditions:, actions:, inclusion: [], exclusion: [])
        @name = name
        @conditions = conditions
        @actions = actions
        @inclusion = inclusion
        @exclusion = exclusion
        validate_exclusive_options!
      end

      private

      def validate_exclusive_options!
        return if (inclusion || []).empty? || (exclusion || []).empty?

        raise Architect::Domain::Errors::InvalidRuleDefinition,
              'exclusion and inclusion are mutually exclusive options'
      end
    end
  end
end
