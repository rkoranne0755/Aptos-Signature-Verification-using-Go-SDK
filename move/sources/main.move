module Admin::verify_signature {

    use std::ed25519;
    use std::bcs;
    use aptos_std::hash;

    const EINVALID_SIGNATURE: u64 = 1001;
    const EALREADY_PRESENT: u64 = 1002;

    const ADMIN_PUBLIC_KEY: vector<u8> = x"d0648a05a4b53cc16d01215acc13ac924777bc58cd8f12919fe72cd7377b5f09";

    // Struct to hold multiple admin public keys
    struct AdminSigner has key {
        signer_key: vector<u8> // List of admin public keys
    }

    struct Message has copy, drop {
        from: address,
        amount: u64,
        nonce: u8
    }

    // List the PubKey

    #[view]
    public fun get_pub_key(): vector<u8> acquires AdminSigner {
        borrow_global<AdminSigner>(@Admin).signer_key
    }

    // Initialize an AdminSigner with a list of public keys
    fun init_module(admin: &signer) {
        move_to(admin, AdminSigner { signer_key: ADMIN_PUBLIC_KEY })
    }

    // Function to verify multiple signatures
    #[view]
    public fun verify_signatures(
        from: address, amount: u64, signature: vector<u8>
    ): bool acquires AdminSigner {
        let msg = Message { from: from, amount: amount, nonce: 1 };

        let msg_bytes = bcs::to_bytes<Message>(&msg);
        let msg_hash = hash::sha2_256(msg_bytes);

        verify(msg_hash, signature)

    }

    fun verify(message_hash: vector<u8>, signature: vector<u8>): bool acquires AdminSigner {
        let signature_ed = ed25519::new_signature_from_bytes(signature);

        // Convert public key bytes to Ed25519 public key
        let pub_key = get_pub_key();
        let public_key_ed = ed25519::new_unvalidated_public_key_from_bytes(pub_key);

        // Verify the signature using the message hash
        let is_valid =
            ed25519::signature_verify_strict(&signature_ed, &public_key_ed, message_hash);

        return is_valid
    }

    public entry fun update_key(_admin: &signer, pub_key: vector<u8>) acquires AdminSigner {
        let admin_key = borrow_global_mut<AdminSigner>(@Admin);
        admin_key.signer_key = pub_key;
    }

    #[test(acc = @Admin)]
    public fun test(acc: &signer) acquires AdminSigner {
        init_module(acc);
        let pub_keys: vector<u8> = vector<u8>[0x01, 0x02, 0x03];
        update_key(acc, pub_keys);
        assert!(pub_keys == get_pub_key(),1001); 
    }
}
