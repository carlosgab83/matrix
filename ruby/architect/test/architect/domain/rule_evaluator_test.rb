# typed: false
# frozen_string_literal: true

require 'test_helper'
require 'architect/shared/app_logger'

module Architect
  module Domain
    class RuleEvaluatorTest < Minitest::Test
      # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
      def test_evaluates_rules
        rule1 = Rule.new(
          name: 'Rule 1',
          conditions: [
            Condition.new(key: 'type', operator: :equals_to, value: 'type1'),
            Condition.new(key: 'price', operator: :less_than, value: 1000)
          ],
          actions: [
            Action.new(action_type: 'notify', action_value: 'price_lower')
          ]
        )

        rule2 = Rule.new(
          name: 'Rule 2',
          conditions: [
            Condition.new(key: 'type', operator: :equals_to, value: 'type2'),
            Condition.new(key: 'sentiment', operator: :not_equals_to, value: 'good')
          ],
          actions: [
            Action.new(action_type: 'notify', action_value: 'bad_news')
          ]
        )

        rules = [rule1, rule2]
        expected = [rule1, rule2, nil]

        json_objs = [
          {
            'type' => 'type1',
            'price' => 900
          },
          {
            'type' => 'type2',
            'sentiment' => 'bad'
          },
          {
            'type' => 'type2',
            'sentiment' => 'good'
          }
        ]

        json_objs.each_with_index do |obj, i|
          if expected[i].nil?
            assert_nil RuleEvaluator.new(rules).evaluate(obj)
          else
            assert_equal RuleEvaluator.new(rules).evaluate(obj), expected[i]
          end
        end
      end
      # rubocop:enable Metrics/AbcSize, Metrics/MethodLength

      # rubocop:disable Metrics/AbcSize, Metrics/MethodLength
      def test_evaluates_rules_with_inclusion_exclusion
        rule1 = Rule.new(
          name: 'Rule 1',
          conditions: [
            Condition.new(key: 'type', operator: :equals_to, value: 'type1'),
            Condition.new(key: 'price', operator: :less_than, value: 1000)
          ],
          actions: [
            Action.new(action_type: 'notify', action_value: 'price_lower')
          ],
          inclusion: ['BTCUSD']
        )

        rule2 = Rule.new(
          name: 'Rule 2',
          conditions: [
            Condition.new(key: 'type', operator: :equals_to, value: 'type2'),
            Condition.new(key: 'sentiment', operator: :not_equals_to, value: 'good')
          ],
          actions: [
            Action.new(action_type: 'notify', action_value: 'bad_news')
          ],
          exclusion: ['BTCUSD']
        )

        rule3 = Rule.new(
          name: 'Rule 3',
          conditions: [
            Condition.new(key: 'type', operator: :equals_to, value: 'type2'),
            Condition.new(key: 'sentiment', operator: :not_equals_to, value: 'good')
          ],
          actions: [
            Action.new(action_type: 'notify', action_value: 'bad_news')
          ],
          inclusion: ['BTCUSD_1']
        )

        rules = [rule1, rule2, rule3]
        expected = [rule1, nil, nil]

        json_objs = [
          {
            'type' => 'type1',
            'price' => 900,
            'symbol' => 'BTCUSD'
          },
          {
            'type' => 'type2',
            'sentiment' => 'bad',
            'symbol' => 'BTCUSD'
          },
          {
            'type' => 'type2',
            'sentiment' => 'good',
            'symbol' => 'BTCUSD_1'
          }
        ]

        json_objs.each_with_index do |obj, i|
          if expected[i].nil?
            assert_nil RuleEvaluator.new(rules).evaluate(obj)
          else
            assert_equal RuleEvaluator.new(rules).evaluate(obj), expected[i]
          end
        end
      end
      # rubocop:enable Metrics/AbcSize, Metrics/MethodLength
    end
  end
end
