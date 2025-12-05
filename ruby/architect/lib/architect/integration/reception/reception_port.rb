# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'

module Architect
  module Integration
    module Reception
      module ReceptionPort
        extend T::Sig
        extend T::Helpers

        interface!

        sig { abstract.params(_block: T.proc.params(message: T.untyped).void).void }
        def receive(&); end
      end

      module ReceptionPortFactory
        extend T::Sig

        sig { params(logger: Logger).returns(Architect::Integration::Reception::KafkaAdapter) }
        def self.new_adapter(logger:)
          KafkaAdapter.new(logger: logger)
        end
      end
    end
  end
end

require 'architect/integration/reception/kafka_adapter'
