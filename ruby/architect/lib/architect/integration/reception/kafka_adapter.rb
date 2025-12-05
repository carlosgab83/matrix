# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'kafka'

module Architect
  module Integration
    module Reception
      class KafkaAdapter
        extend T::Sig
        include ReceptionPort

        sig { params(logger: Logger).void }
        def initialize(logger:)
          @logger = T.let(logger, Logger)
          @kafka = T.let(::Kafka.new(ENV.fetch('KAFKA_BOOTSTRAP_SERVERS').split(',')), ::Kafka::Client)
          @topic = T.let('snapshot.new', String)
          @group_id = T.let('snapshot-processor', String)
          @logger.debug("Kafka adapter initialized for topic: #{@topic}")
        end

        sig { override.params(_block: T.proc.params(message: T.untyped).void).void }
        def receive(&)
          consumer = @kafka.consumer(group_id: @group_id)
          consumer.subscribe(@topic)
          @logger.info('Waiting for messages...')

          # This blocks forever, processing messages as they arrive
          consumer.each_message do |message|
            @logger.debug("Received message: partition=#{message.partition}, offset=#{message.offset}")
            begin
              yield(message)
            rescue StandardError => e
              @logger.error("Error processing message: #{e.message}")
            end
          end
        rescue StandardError => e
          @logger.fatal("FATAL ERROR: #{e.class} - #{e.message}")
        ensure
          consumer.stop if consumer
          @logger.info('Consumer stopped')
        end
      end
    end
  end
end
