# rubocop:disable all
# typed: false

require 'architect/domain/dsl'

include Architect::Domain::DSL

rule 'fixture_rule_2' do
  for_type :some_type2
  condition :volume, :greater_than, 500
  except [:MSFT]
end
