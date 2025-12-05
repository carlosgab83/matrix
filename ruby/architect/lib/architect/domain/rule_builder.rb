# typed: false
# frozen_string_literal: true

require 'architect/domain/condition'
require 'architect/domain/action'
require 'architect/domain/rule'

module Architect
  module Domain
    class RuleBuilder
      attr_reader :name

      def initialize(name)
        @name = name
        @conditions = []
        @actions = []
      end

      def for_type(type)
        @conditions << Condition.new(key: :type, operator: :equals_to, value: type.to_sym)
      end

      def condition(key, operator, value)
        @conditions << Condition.new(key: key.to_sym, operator: operator.to_sym, value: value)
      end

      def only(stock_symbol)
        @inclusion = stock_symbol.is_a?(Array) ? stock_symbol : [stock_symbol]
      end

      def except(stock_symbol)
        @exclusion = stock_symbol.is_a?(Array) ? stock_symbol : [stock_symbol]
      end

      def method_missing(method, *args)
        case method.to_sym
        when :action
          @actions << Action.new(action_type: args[0].to_sym, action_value: args[1])
        end
      end

      def respond_to_missing?(method, _include_private = false)
        [:action].include?(method.to_sym)
      end

      def build
        Rule.new(name: name, conditions: @conditions, actions: @actions, inclusion: @inclusion, exclusion: @exclusion)
      end
    end
  end
end
