# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/errors/invalid_rule_definition'

module Architect
  module Domain
    class Rule
      extend T::Sig

      sig { returns(String) }
      attr_reader :name

      sig { returns(T::Array[Condition]) }
      attr_reader :conditions

      sig { returns(T::Array[Action]) }
      attr_reader :actions

      sig { returns(T::Array[String]) }
      attr_reader :inclusion

      sig { returns(T::Array[String]) }
      attr_reader :exclusion

      sig do
        params(name: String, conditions: T::Array[Condition], actions: T::Array[Action], inclusion: T::Array[String],
               exclusion: T::Array[String]).void
      end
      def initialize(name:, conditions:, actions:, inclusion: [], exclusion: [])
        @name = T.let(name, String)
        @conditions = T.let(conditions, T::Array[Condition])
        @actions = T.let(actions, T::Array[Action])
        @inclusion = T.let(inclusion, T::Array[String])
        @exclusion = T.let(exclusion, T::Array[String])
        validate_exclusive_options!
      end

      private

      sig { void }
      def validate_exclusive_options!
        return if inclusion.empty? || exclusion.empty?

        raise Architect::Domain::Errors::InvalidRuleDefinition,
              'exclusion and inclusion are mutually exclusive options'
      end
    end
  end
end
