#ifndef EASYQ_BRIDGE_H
#define EASYQ_BRIDGE_H

#ifdef __cplusplus
extern "C" {
#endif

#include <stdint.h>

// Status codes
#define EASYQ_SUCCESS 0
#define EASYQ_ERROR_GENERAL 1
#define EASYQ_ERROR_NOT_INITIALIZED 2
#define EASYQ_ERROR_INVALID_ARGUMENT 3
#define EASYQ_ERROR_RUNTIME 4
#define EASYQ_ERROR_TIMEOUT 5
#define EASYQ_ERROR_AUTHENTICATION 6
#define EASYQ_ERROR_CONNECTION 7

// Initialize the EasyQ runtime
int EasyQ_Initialize();

// Shutdown the EasyQ runtime
void EasyQ_Shutdown();

// Configure the connection to a quantum computing resource
int EasyQ_ConfigureConnection(const char* config_json);

// Generate a random integer between min and max (inclusive)
int EasyQ_GenerateRandomInt(int min, int max, int* result);

// Generate random bytes
int EasyQ_GenerateRandomBytes(int length, unsigned char* buffer);

// Perform a quantum search operation
int EasyQ_Search(const char* items_json, const char* predicate_json, 
                 const char* options_json, char** result_json);

// Generate a cryptographically secure key using quantum key distribution
int EasyQ_GenerateKey(const char* options_json, char** result_json);

// Verify the security of a quantum channel
int EasyQ_VerifyChannelSecurity(const char* options_json, char** result_json);

// Free a string allocated by the EasyQ runtime
void EasyQ_FreeString(char* str);

#ifdef __cplusplus
}
#endif

#endif // EASYQ_BRIDGE_H