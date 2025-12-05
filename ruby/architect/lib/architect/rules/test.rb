# rubocop:disable all
# typed: false
# frozen_string_literal: true

require 'architect/domain/dsl'

include Architect::Domain::DSL

rule 'Testing Rule' do
  for_type 'test'
  condition 'price', 'greater_than', 1
  only %w[BTCUSD]
  action 'notify', 'test'
end
