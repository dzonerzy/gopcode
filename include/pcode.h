#ifndef PCODE_NATIVE_H
#define PCODE_NATIVE_H

#include <stdint.h>
#include <stddef.h>

#ifdef __cplusplus
extern "C"
{
#endif

    typedef struct PcodeContext PcodeContext;
    typedef struct NativeAddrSpace NativeAddrSpace;

    typedef enum CAddrSpaceFlags
    {
        big_endian = 1,             ///< Space is big endian if set, little endian otherwise
        heritaged = 2,              ///< This space is heritaged
        does_deadcode = 4,          ///< Dead-code analysis is done on this space
        programspecific = 8,        ///< Space is specific to a particular loadimage
        reverse_justification = 16, ///< Justification within aligned word is opposite of endianness
        formal_stackspace = 0x20,   ///< Space attached to the formal \b stack \b pointer
        overlay = 0x40,             ///< This space is an overlay of another space
        overlaybase = 0x80,         ///< This is the base space for overlay space(s)
        truncated = 0x100,          ///< Space is truncated from its original size, expect pointers larger than this size
        hasphysical = 0x200,        ///< Has physical memory associated with it
        is_otherspace = 0x400,      ///< Quick check for the OtherSpace derived class
        has_nearpointers = 0x800    ///< Does there exist near pointers into this space
    } CAddrSpaceFlags;

    typedef enum CTranslateFlags
    {
        bb_terminating = 1,
    } CTranslateFlags;

    typedef struct
    {
        const char *name;
        uint32_t index;
        uint32_t address_size;
        uint32_t word_size;
        uint32_t flags;
        uint64_t highest;
        uint64_t pointer_lower_bound;
        uint64_t pointer_upper_bound;
        NativeAddrSpace *n_space;
    } AddrSpaceC;

    typedef struct
    {
        AddrSpaceC *space;
        unsigned long long offset;
        int32_t size;
    } VarnodeDataC;

    typedef struct
    {
        uint32_t opcode;
        VarnodeDataC *output;
        VarnodeDataC *inputs;
        uint32_t num_inputs;
    } PcodeOpC;

    typedef struct
    {
        PcodeOpC *ops;
        uint32_t num_ops;
    } PcodeTranslationC;

    typedef struct
    {
        uint64_t address;
        uint64_t length;
        const char *mnemonic;
        const char *body;
    } DisassemblyInstructionC;

    typedef struct
    {
        DisassemblyInstructionC *instructions;
        uint32_t num_instructions;
    } PcodeDisassemblyC;

    typedef struct RegisterInfoC
    {
        VarnodeDataC varnode;
        const char *name;
    } RegisterInfoC;

    typedef struct RegisterInfoListC
    {
        RegisterInfoC *registers;
        uint32_t count;
    } RegisterInfoListC;

    // Initialize and release the context
    PcodeContext *pcode_context_create(unsigned char *slaBytes, size_t slaSize);
    void pcode_context_free(PcodeContext *ctx);
    void pcode_context_set_variable_default(PcodeContext *ctx, const char *nm, uint32_t val);
    RegisterInfoListC *pcode_context_get_all_registers(PcodeContext *ctx);
    const char *pcode_context_get_register_name(PcodeContext *ctx, NativeAddrSpace *space, unsigned long long offset, int32_t size);

    // Disassemble code
    PcodeDisassemblyC *pcode_disassemble(PcodeContext *ctx, const char *bytes, unsigned int num_bytes, unsigned long long base_address, unsigned int max_instructions);
    void pcode_disassembly_free(PcodeDisassemblyC *disas);

    // Translate code
    PcodeTranslationC *pcode_translate(PcodeContext *ctx, const char *bytes, unsigned int num_bytes, unsigned long long base_address, unsigned int max_instructions, uint32_t flags);
    void pcode_translation_free(PcodeTranslationC *trans);

    // VarNode code
    const char *pcode_varcode_get_register_name(NativeAddrSpace *space, unsigned long long offset, int32_t size);
    AddrSpaceC *pcode_varnode_get_space_from_const(unsigned long long offset);

#ifdef __cplusplus
}
#endif

#endif // PCODE_NATIVE_H
