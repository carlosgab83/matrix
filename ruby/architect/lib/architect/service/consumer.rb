# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/domain/registry'
require 'architect/domain/rule_evaluator'
require 'architect/integration/publication/publication_port'

module Architect
  module Service
    class Consumer
      extend T::Sig

      sig do
        params(registry: T.class_of(Architect::Domain::Registry), logger: Logger,
               receptor: Architect::Integration::Reception::ReceptionPort, publisher: Architect::Integration::Publication::PublicationPort).void
      end
      def initialize(registry:, logger:, receptor:, publisher:)
        @registry = T.let(registry, T.class_of(Architect::Domain::Registry))
        @receptor = T.let(receptor, Architect::Integration::Reception::ReceptionPort)
        @publisher = T.let(publisher, Architect::Integration::Publication::PublicationPort)
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
            @logger.info("Rule '#{rule.name}' matched for snapshot")
            build_and_publish_notification(data, rule)
          else
            @logger.info('No rules matched, ignoring snapshot')
          end

          @logger.debug('Message processed successfully')
        end
      end

      private

      sig { params(snapshot: T::Hash[String, T.untyped], rule: Architect::Domain::Rule).void }
      def build_and_publish_notification(snapshot, rule)
        notification = {
          'symbol' => snapshot['symbol'],
          'price' => snapshot['price'],
          'timestamp' => snapshot['timestamp'],
          'actions' => rule.actions.map do |action|
            { 'type' => action.action_type.to_s, 'value' => action.action_value }
          end
        }

        @publisher.publish(notification)
        @logger.info("Notification published for rule '#{rule.name}'")
      rescue StandardError => e
        @logger.error("Failed to publish notification: #{e.message}")
      end
    end
  end
end
