# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/condition'
require 'architect/domain/action'
require 'architect/domain/rule'

module Architect
  module Domain
    class RuleBuilder
      extend T::Sig

      sig { returns(String) }
      attr_reader :name

      sig { params(name: String).void }
      def initialize(name)
        @name = T.let(name, String)
        @conditions = T.let([], T::Array[Condition])
        @actions = T.let([], T::Array[Action])
        @inclusion = T.let([], T::Array[String])
        @exclusion = T.let([], T::Array[String])
      end

      sig { params(type: T.any(String, Symbol)).void }
      def for_type(type)
        @conditions << Condition.new(key: 'type', operator: :equals_to, value: type.to_s)
      end

      sig { params(key: T.any(String, Symbol), operator: T.any(String, Symbol), value: T.untyped).void }
      def condition(key, operator, value)
        @conditions << Condition.new(key: key.to_s, operator: operator.to_sym, value: value)
      end

      sig { params(stock_symbol: T.any(String, T::Array[String])).void }
      def only(stock_symbol)
        @inclusion = stock_symbol.is_a?(Array) ? stock_symbol : [stock_symbol]
      end

      sig { params(stock_symbol: T.any(String, T::Array[String])).void }
      def except(stock_symbol)
        @exclusion = stock_symbol.is_a?(Array) ? stock_symbol : [stock_symbol]
      end

      sig { params(method: Symbol, args: T.untyped).void }
      def method_missing(method, *args)
        case method.to_sym
        when :action
          @actions << Action.new(action_type: args[0].to_s, action_value: args[1].to_s)
        end
      end

      sig { params(method: T.any(Symbol, String), _include_private: T::Boolean).returns(T::Boolean) }
      def respond_to_missing?(method, _include_private = false)
        [:action].include?(method.to_sym)
      end

      sig { returns(Rule) }
      def build
        Rule.new(name: name, conditions: @conditions, actions: @actions, inclusion: @inclusion, exclusion: @exclusion)
      end
    end
  end
end
