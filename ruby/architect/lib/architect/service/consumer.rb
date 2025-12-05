# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/registry'
require 'architect/domain/rule_evaluator'

module Architect
  module Service
    class Consumer
      extend T::Sig

      sig { params(registry: T.class_of(Architect::Domain::Registry), logger: Logger, receptor: Architect::Integration::Reception::ReceptionPort, publisher: T.untyped).void }
      def initialize(registry:, logger:, receptor:, publisher:)
        @registry = T.let(registry, T.class_of(Architect::Domain::Registry))
        @receptor = T.let(receptor, Architect::Integration::Reception::ReceptionPort)
        @publisher = T.let(publisher, T.untyped)
        @logger = T.let(logger, Logger)
      end

      sig { void }
      def call
        @logger.info('Starting message consumer...')
        @receptor.receive do |msg|
          @logger.debug('Processing incoming message...')
          data = JSON.parse(msg.value)
          @logger.info("Received snapshot: #{data.inspect}")
          rule = Architect::Domain::RuleEvaluator.new(@registry.rules).evaluate(data)
          if rule
            @logger.info("Processing #{data} with rule #{rule.name}")
            # TODO: Process and send to tank
          else
            @logger.info("Ignoring #{data}")
          end

          @logger.debug('Message processed successfully')
        end
      end
    end
  end
end
