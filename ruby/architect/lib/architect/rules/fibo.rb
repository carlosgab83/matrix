# rubocop:disable all
# typed: false
# frozen_string_literal: true

require 'architect/domain/dsl'

include Architect::Domain::DSL

rule 'Fibonacci Extension' do
  for_type 'fibonacci_extension'
  condition 'value', 'equals_to', 0.618
  only %w[BTCUSD aa]
  action 'notify', 'fibonaccy_extension'
end

rule 'Fibonacci Retracement' do
end
