# typed: false
# frozen_string_literal: true

require 'architect/domain/errors/invalid_rule_definition'

module Architect
  module Domain
    class Registry
      class << self
        @unique_validator = {}

        def rules
          @rules ||= []
        end

        def add(rule)
          Architect::Shared::AppLogger.instance.info("Adding rule: #{rule.name}")
          if @unique_validator[rule.name]
            raise Architect::Domain::Errors::InvalidRuleDefinition,
                  "Duplicted rule name #{rule.name}"
          end

          @unique_validator[rule.name] = true
          rules << rule
        end

        def clear!
          @unique_validator = {}
          @rules = []
        end

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
