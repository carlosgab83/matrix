# typed: strict
# frozen_string_literal: true

require 'sorbet-runtime'

module Architect
  module Domain
    class Action
      extend T::Sig

      sig { returns(String) }
      attr_reader :action_type

      sig { returns(String) }
      attr_reader :action_value

      sig { params(action_type: String, action_value: String).void }
      def initialize(action_type:, action_value:)
        @action_type = T.let(action_type, String)
        @action_value = T.let(action_value, String)
      end
    end
  end
end
