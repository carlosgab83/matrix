# typed: false
# frozen_string_literal: true

require 'architect/domain/condition'
require 'architect/domain/action'
require 'architect/domain/rule'

module Architect
  module Domain
    class RuleEvaluator
      attr_reader :rules

      def initialize(rules)
        @rules = rules
      end

      def evaluate(snapshot)
        rules.each do |rule|
          return rule if matches?(rule, snapshot)
        end

        nil
      end

      # rubocop:disable Metrics/AbcSize, Metrics/CyclomaticComplexity, Metrics/PerceivedComplexity
      def matches?(rule, snapshot)
        rule.conditions.each do |condition|
          return false unless snapshot[condition.key]
          return false if rule.inclusion&.any? && !rule.inclusion.include?(snapshot['symbol'])
          return false if rule.exclusion&.any? && rule.exclusion.include?(snapshot['symbol'])

          case condition.operator
          when :equals_to
            return false if snapshot[condition.key] != condition.value
          when :not_equals_to
            return false if snapshot[condition.key] == condition.value
          when :greater_than
            return false if snapshot[condition.key] <= condition.value
          when :greater_or_equal
            return false if snapshot[condition.key] < condition.value
          when :less_than
            return false if snapshot[condition.key] >= condition.value
          when :less_or_equal
            return false if snapshot[condition.key] > condition.value
          end
        end

        true
      end
      # rubocop:enable Metrics/AbcSize, Metrics/CyclomaticComplexity, Metrics/PerceivedComplexity
    end
  end
end
