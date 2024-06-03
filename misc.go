package radir

func isDigit(r rune) bool {                                             
        return '0' <= r && r <= '9'                                     
}                                                                       
                                                                        
func isAlpha(r rune) bool {                                             
        return isLAlpha(r) || isUAlpha(r)                               
}                                                                       
                                                                        
func isUAlpha(r rune) bool {                                            
        return 'A' <= r && r <= 'Z'                                     
}                                                                       
                                                                        
func isLAlpha(r rune) bool {                                            
        return 'a' <= r && r <= 'z'                                     
}

/*                                                                      
oidFromValue returns a string following an attempt to associate a numeric
OID based on the input name and section. A zero string indicates no match.
*/                                                                      
func oidFromValue(name, section string, v []string) (oid string) {      
        if len(section) == 0 {                                          
                return                                                  
        }                                                               
                                                                        
        switch len(v) {                                                 
        case 0:                                                         
                // no name to match                                     
        case 1:                                                         
                if eq(v[0], name) {                                     
                        oid = PrefixOID + `.` + section                 
                }                                                       
        default:                                                        
                for _, val := range v {                                 
                        if eq(val, name) {                              
                                oid = PrefixOID + `.` + section         
                                break                                   
                        }                                               
                }                                                       
        }                                                               
                                                                        
        return                                                          
}                                                                       
                                                                        
/*                                                                      
defByOID returns the principal name of the specified schema definition OID.
                                                                        
Valid type (typ) values are `at`, `oc`, and `nf`.                       
*/                                                                      
func defByOID(typ, oid string) (name string) {                          
        if _, ok := schema[typ]; !ok {                                  
                return                                                  
        }                                                               
                                                                        
        if val, ok := schema[typ][oid]; ok {                            
                if len(val) > 0 {                                       
                        name = val[0]                                   
                }                                                       
        }                                                               
                                                                        
        return                                                          
}

/*                                                                      
defByName returns a string (oid) associated with the input name.        
*/                                                                      
func defByName(typ, name string) (oid string) {                         
        if _, ok := schema[typ]; !ok || len(name) == 0 {                
                return                                                  
        }                                                               
                                                                        
        for section, v := range schema[typ] {                           
                if oid = oidFromValue(name, section, v); len(oid) > 0 { 
                        break                                           
                }                                                       
        }                                                               
                                                                        
        return                                                          
}

