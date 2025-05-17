#ifndef EASYQ_BRIDGE_H
#define EASYQ_BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

/* Basic functions */
int EasyQ_Initialize();
void EasyQ_Shutdown();
int EasyQ_ConfigureConnection(const char* config_json);
void EasyQ_FreeString(char* str);

/* Quantum Search */
int EasyQ_Search(
    const char* items_json, 
    const char* predicate_json,
    const char* options_json, 
    char** result_json
);

/* Quantum Random Number Generation */
int EasyQ_GenerateRandomInt(int min, int max, int* result);
int EasyQ_GenerateRandomBytes(int length, unsigned char* buffer);

/* Quantum Key Distribution */
int EasyQ_GenerateKey(
    const char* options_json, 
    char** result_json
);

/* Error codes */
#define EASYQ_SUCCESS 0
#define EASYQ_ERROR_GENERAL 1
#define EASYQ_ERROR_NOT_INITIALIZED 2
#define EASYQ_ERROR_INVALID_ARGUMENT 3
#define EASYQ_ERROR_RUNTIME 4
#define EASYQ_ERROR_TIMEOUT 5

#ifdef __cplusplus
}
#endif

#endif /* EASYQ_BRIDGE_H */