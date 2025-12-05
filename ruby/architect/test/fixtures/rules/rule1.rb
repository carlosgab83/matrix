# rubocop:disable all
# typed: false

require 'architect/domain/dsl'

include Architect::Domain::DSL

rule 'fixture_rule_1' do
  for_type :some_type1
  condition :price, :greater_than, 50
  only [:AAPL]
end
