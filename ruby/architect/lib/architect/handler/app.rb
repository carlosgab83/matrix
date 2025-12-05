# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'
require 'architect/integration/reception/reception_port'
require 'architect/service/consumer'
require 'architect/shared/app_logger'

module Architect
  module Handler
    class App
      extend T::Sig

      sig { void }
      def self.run
        logger = Architect::Shared::AppLogger.instance
        logger.info('Initializing Architect Handler')

        receptor = Architect::Integration::Reception::ReceptionPortFactory.new_adapter(logger: logger)

        Architect::Service::Consumer.new(
          logger: logger,
          receptor: receptor,
          publisher: nil
        ).call
      end
    end
  end
end
