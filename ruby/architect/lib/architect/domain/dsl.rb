# typed: false
# frozen_string_literal: true

require "architect/domain/rule_builder"

module Architect
  module Domain
    module DSL
      def rule(name, &block)
        builder = RuleBuilder.new(name)
        builder.instance_eval(&block)
        Registry.add(builder.build)



      end
    end
  end
end
