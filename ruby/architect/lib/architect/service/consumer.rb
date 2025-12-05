# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'

module Architect
  module Service
    class Consumer
      extend T::Sig

      sig { params(logger: Logger, receptor: Architect::Integration::Reception::ReceptionPort, publisher: T.untyped).void }
      def initialize(logger:, receptor:, publisher:)
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
          # TODO: Process and publish to tank
          @logger.debug('Message processed successfully')
        end
      end
    end
  end
end
