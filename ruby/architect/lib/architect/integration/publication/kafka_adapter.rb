# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'kafka'
require 'json'

module Architect
  module Integration
    module Publication
      class KafkaAdapter
        extend T::Sig
        include PublicationPort

        NOTIFICATION_TOPIC = 'notification.new'

        sig { params(logger: Logger).void }
        def initialize(logger:)
          @logger = T.let(logger, Logger)
          @kafka = T.let(Kafka.new(['nabucodonosor:9092'], client_id: 'architect-publisher'), Kafka::Client)
          @producer = T.let(@kafka.producer, Kafka::Producer)
        end

        sig { override.params(payload: T::Hash[String, T.untyped]).void }
        def publish(payload)
          message = JSON.generate(payload)
          @producer.produce(message, topic: NOTIFICATION_TOPIC)
          @producer.deliver_messages
          @logger.info("Published message to '#{NOTIFICATION_TOPIC}': #{payload.inspect}")
        rescue StandardError => e
          @logger.error("Error publishing message: #{e.message}")
          raise
        end

        sig { void }
        def close
          @producer.shutdown
        end
      end
    end
  end
end
