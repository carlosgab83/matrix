# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/rule_builder'

module Architect
  module Domain
    module DSL
      extend T::Sig

      sig { params(name: String, block: T.proc.bind(RuleBuilder).void).void }
      def rule(name, &block)
        builder = RuleBuilder.new(name)
        builder.instance_eval(&block)
        Registry.add(builder.build)




      end
    end
  end
end
