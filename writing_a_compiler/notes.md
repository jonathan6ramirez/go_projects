# Writing a Compiler Notes

---

## Chapter: Functions in the VM:

>[!QUESTION]
> What are we going to do to implement calling the calling of functions
> and the execution of them.

>[!ANSWER]
> So the way to look at it is like frames. There is the 'mainframe' and then there
> is the rest of the 'frames'. The rest of these frames are the instructions to their
> bodies. So all we have to do is change the 'instructions' or 'frame' of the VM and the
> 'pointer' or the 'ip' variable.


>[!QUESTION]
> So how do we handle calling functions?

>[!ANSWER]
> We handle function calls by checking for 'OpReturnValue' inside our Run() function
> then we get the value by popping off the top of the stack and pop the functions frame of the frame stack.
> Lastly we pop again to get the 'CompiledFunction' off the stack.

>[!QUESTION]
> So how do we handle functions that have no return value?
```monkey (made up language)
let noReturn = fn() {};
noReturn();

Null
```
>[!ANSWER]
> So we do that by checking for the Return opcode inside our run function
> then we pop the current frame and pop the CompiledFunction in the 'mainframe'

>[!QUESTION]
> How do we handle different scopes in according to how the functions are nested
> and where they are called?

>[!ANSWER]
> We handle them by checking the scope that is assigned to the 'LetStatement'
> and the 'Identifier' nodes. If the scope is pointing to local then we emit
> a local scope and vice versa.
