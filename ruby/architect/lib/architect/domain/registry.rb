# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/errors/invalid_rule_definition'

module Architect
  module Domain
    class Registry
      class << self
        extend T::Sig

        sig { returns(T::Hash[String, T::Boolean]) }
        def unique_validator
          @unique_validator = T.let(@unique_validator, T.nilable(T::Hash[String, T::Boolean]))
          @unique_validator ||= {}
        end

        sig { returns(T::Array[Rule]) }
        def rules
          @rules = T.let(@rules, T.nilable(T::Array[Rule]))
          @rules ||= []
        end

        sig { params(rule: Rule).void }
        def add(rule)
          Architect::Shared::AppLogger.instance.info("Adding rule: #{rule.name}")
          if unique_validator[rule.name]
            raise Architect::Domain::Errors::InvalidRuleDefinition,
                  "Duplicted rule name #{rule.name}"
          end

          unique_validator[rule.name] = true
          rules << rule
        end

        sig { void }
        def clear!
          @unique_validator = {}
          @rules = []
        end

        sig { params(path: String).void }
        def reload_from_dir!(path)
          clear!
          Dir["#{path}/**/*.rb"].each do |file|
            load file
          end
        end
      end
    end
  end
end
