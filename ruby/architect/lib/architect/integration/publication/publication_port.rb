# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'

module Architect
  module Integration
    module Publication
      module PublicationPort
        extend T::Sig
        extend T::Helpers

        interface!

        sig { abstract.params(payload: T::Hash[String, T.untyped]).void }
        def publish(payload); end
      end

      module PublicationPortFactory
        extend T::Sig

        sig { params(logger: Logger).returns(Architect::Integration::Publication::KafkaAdapter) }
        def self.new_adapter(logger:)
          KafkaAdapter.new(logger: logger)
        end
      end
    end
  end
end

require 'architect/integration/publication/kafka_adapter'
