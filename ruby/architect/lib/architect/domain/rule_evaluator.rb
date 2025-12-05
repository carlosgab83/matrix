# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/condition'
require 'architect/domain/action'
require 'architect/domain/rule'

module Architect
  module Domain
    class RuleEvaluator
      extend T::Sig

      sig { returns(T::Array[Rule]) }
      attr_reader :rules

      sig { params(rules: T::Array[Rule]).void }
      def initialize(rules)
        @rules = T.let(rules, T::Array[Rule])
      end

      sig { params(snapshot: T::Hash[String, T.untyped]).returns(T.any(Rule, NilClass)) }
      def evaluate(snapshot)
        rules.each do |rule|
          return rule if matches?(rule, snapshot)
        end

        nil
      end

      # rubocop:disable Metrics/AbcSize, Metrics/CyclomaticComplexity, Metrics/PerceivedComplexity
      sig { params(rule: Rule, snapshot: T.untyped).returns(T::Boolean) }
      def matches?(rule, snapshot)
        rule.conditions.each do |condition|
          return false unless snapshot[condition.key]
          return false if rule.inclusion.any? && !rule.inclusion.include?(snapshot['symbol'])
          return false if rule.exclusion.any? && rule.exclusion.include?(snapshot['symbol'])

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
